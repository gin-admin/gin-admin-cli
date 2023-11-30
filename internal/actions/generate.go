package actions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/parser"
	"github.com/gin-admin/gin-admin-cli/v10/internal/schema"
	"github.com/gin-admin/gin-admin-cli/v10/internal/tfs"
	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
	"go.uber.org/zap"
)

type GenerateConfig struct {
	Dir         string
	TplType     string
	ModuleName  string
	ModulePath  string
	WirePath    string
	SwaggerPath string
	FEDir       string
	FETplType   string
}

func Generate(cfg GenerateConfig) *GenerateAction {
	return &GenerateAction{
		logger:           zap.S().Named("[GEN]"),
		cfg:              &cfg,
		fs:               tfs.Ins,
		rootImportPath:   parser.GetRootImportPath(cfg.Dir),
		moduleImportPath: parser.GetModuleImportPath(cfg.Dir, cfg.ModulePath, cfg.ModuleName),
		UtilImportPath:   parser.GetUtilImportPath(cfg.Dir, cfg.ModulePath),
	}
}

type GenerateAction struct {
	logger           *zap.SugaredLogger
	cfg              *GenerateConfig
	fs               tfs.FS
	rootImportPath   string
	moduleImportPath string
	UtilImportPath   string
}

// Run generate command
func (a *GenerateAction) RunWithConfig(ctx context.Context, cfgName string) error {
	var parseFile = func(name string) error {
		var data []*schema.S
		switch filepath.Ext(name) {
		case ".json":
			if err := utils.ParseJSONFile(name, &data); err != nil {
				return err
			}
			return a.run(ctx, data)
		case ".yaml", "yml":
			if err := utils.ParseYAMLFile(name, &data); err != nil {
				return err
			}
			return a.run(ctx, data)
		default:
			a.logger.Warnf("Ignore file %s, only support json/yaml/yml", name)
		}
		return nil
	}

	if utils.IsDir(cfgName) {
		return filepath.WalkDir(cfgName, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			return parseFile(path)
		})
	}

	return parseFile(cfgName)
}

func (a *GenerateAction) RunWithStruct(ctx context.Context, structName, comment, output string) error {
	var outputs []string
	if output != "" {
		outputs = strings.Split(output, ",")
	}

	input := []*schema.S{
		{Name: structName, Comment: comment, Outputs: outputs},
	}

	return a.run(ctx, input)
}

func (a *GenerateAction) run(ctx context.Context, data []*schema.S) error {
	for _, d := range data {
		if err := a.generate(ctx, d); err != nil {
			return err
		}

		if d.GenerateFE {
			if err := a.generateFE(ctx, d); err != nil {
				return err
			}
		}
	}

	modsTplData, err := parser.ModifyModsFile(ctx, parser.BasicArgs{
		Dir:        a.cfg.Dir,
		ModuleName: a.cfg.ModuleName,
		ModulePath: a.cfg.ModulePath,
		Flag:       parser.AstFlagGen,
	})
	if err != nil {
		a.logger.Errorf("Failed to modify mods file, err: %s", err)
		return err
	}

	err = a.write(ctx, a.cfg.ModuleName, "", parser.FileForMods, modsTplData, false)
	if err != nil {
		return err
	}

	return a.execWireAndSwag(ctx)
}

func (a *GenerateAction) getGoTplFile(pkgName, tplType string) string {
	pkgName = fmt.Sprintf("%s.go.tpl", pkgName)
	if tplType == "" && a.cfg.TplType != "" {
		tplType = a.cfg.TplType
	}

	if tplType != "" {
		p := filepath.Join(tplType, pkgName)
		if ok, _ := utils.ExistsFile(p); ok {
			return p
		}
		return filepath.Join("default", pkgName)
	}
	return pkgName
}

func (a GenerateAction) getAbsPath(file string) (string, error) {
	modPath := a.cfg.ModulePath
	file = filepath.Join(a.cfg.Dir, modPath, file)
	fullpath, err := filepath.Abs(file)
	if err != nil {
		a.logger.Errorf("Failed to get abs path, err: %s, #file %s", err, file)
		return "", err
	}
	return fullpath, nil
}

func (a *GenerateAction) write(ctx context.Context, moduleName, structName, tpl string, data []byte, checkExists bool) error {
	file, err := parser.ParseFilePathFromTpl(moduleName, structName, tpl)
	if err != nil {
		a.logger.Errorf("Failed to parse file path from tpl, err: %s, #tpl %s", err, tpl)
		return err
	}

	file, err = a.getAbsPath(file)
	if err != nil {
		return err
	}

	exists, err := utils.ExistsFile(file)
	if err != nil {
		return err
	}
	if checkExists && exists {
		a.logger.Infof("File exists, skip, #file %s", file)
		return nil
	}

	a.logger.Infof("Write file: %s", file)
	if !exists {
		err = os.MkdirAll(filepath.Dir(file), os.ModePerm)
		if err != nil {
			a.logger.Errorf("Failed to create dir, err: %s, #dir %s", err, filepath.Dir(file))
			return err
		}
	}

	if exists {
		if err := os.Remove(file); err != nil {
			a.logger.Errorf("Failed to remove file, err: %s, #file %s", err, file)
			return err
		}
	}

	if err := utils.WriteFile(file, data); err != nil {
		a.logger.Errorf("Failed to write file, err: %s, #file %s", err, file)
		return err
	}

	if err := utils.ExecGoFormat(file); err != nil {
		a.logger.Errorf("Failed to exec go format, err: %s, #file %s", err, file)
		return nil
	}

	if err := utils.ExecGoImports(a.cfg.Dir, file); err != nil {
		a.logger.Errorf("Failed to exec go imports, err: %s, #file %s", err, file)
		return nil
	}
	return nil
}

func (a *GenerateAction) generate(ctx context.Context, dataItem *schema.S) error {
	dataItem = dataItem.Format()
	dataItem.RootImportPath = a.rootImportPath
	dataItem.ModuleName = a.cfg.ModuleName
	dataItem.ModuleImportPath = a.moduleImportPath
	dataItem.UtilImportPath = a.UtilImportPath

	genPackages := parser.StructPackages
	if len(dataItem.Outputs) > 0 {
		genPackages = dataItem.Outputs
	}

	for _, pkgName := range genPackages {
		tplName := a.getGoTplFile(pkgName, dataItem.TplType)
		tplData, err := a.fs.ParseTpl(tplName, dataItem)
		if err != nil {
			a.logger.Errorf("Failed to parse tpl, err: %s, #struct %s, #tpl %s", err, dataItem.Name, tplName)
			return err
		}

		err = a.write(ctx, dataItem.ModuleName, dataItem.Name, parser.StructPackageTplPaths[pkgName], tplData, true)
		if err != nil {
			return err
		}
	}

	basicArgs := parser.BasicArgs{
		Dir:         a.cfg.Dir,
		ModuleName:  dataItem.ModuleName,
		ModulePath:  a.cfg.ModulePath,
		StructName:  dataItem.Name,
		GenPackages: genPackages,
		Flag:        parser.AstFlagGen,
	}
	moduleMainTplData, err := parser.ModifyModuleMainFile(ctx, basicArgs)
	if err != nil {
		a.logger.Errorf("Failed to modify module main file, err: %s, #struct %s", err, dataItem.Name)
		return err
	}

	err = a.write(ctx, dataItem.ModuleName, dataItem.Name, parser.FileForModuleMain, moduleMainTplData, false)
	if err != nil {
		return err
	}

	moduleWireTplData, err := parser.ModifyModuleWireFile(ctx, basicArgs)
	if err != nil {
		a.logger.Errorf("Failed to modify module wire file, err: %s, #struct %s", err, dataItem.Name)
		return err
	}

	err = a.write(ctx, dataItem.ModuleName, dataItem.Name, parser.FileForModuleWire, moduleWireTplData, false)
	if err != nil {
		return err
	}

	return nil
}

func (a *GenerateAction) execWireAndSwag(ctx context.Context) error {
	if p := a.cfg.WirePath; p != "" {
		if err := utils.ExecWireGen(a.cfg.Dir, p); err != nil {
			a.logger.Errorf("Failed to exec wire, err: %s, #wirePath %s", err, p)
		}
	}

	if p := a.cfg.SwaggerPath; p != "" {
		if err := utils.ExecSwagGen(a.cfg.Dir, "main.go", p); err != nil {
			a.logger.Errorf("Failed to exec swag, err: %s, #swaggerPath %s", err, p)
		}
	}

	return nil
}

func (a *GenerateAction) generateFE(ctx context.Context, dataItem *schema.S) error {
	for tpl, file := range dataItem.FEMapping {
		tplPath := filepath.Join(a.cfg.FETplType, tpl)
		tplData, err := a.fs.ParseTpl(tplPath, dataItem)
		if err != nil {
			a.logger.Errorf("Failed to parse tpl, err: %s, #struct %s, #tpl %s", err, dataItem.Name, tplPath)
			return err
		}

		file, err := filepath.Abs(filepath.Join(a.cfg.FEDir, file))
		if err != nil {
			return err
		}

		exists, err := utils.ExistsFile(file)
		if err != nil {
			return err
		}
		if exists {
			a.logger.Infof("File exists, skip, #file %s", file)
			continue
		}

		_ = os.MkdirAll(filepath.Dir(file), os.ModePerm)
		if err := utils.WriteFile(file, tplData); err != nil {
			a.logger.Errorf("Failed to write file, err: %s, #file %s", err, file)
			return err
		}
		a.logger.Info("Write file: ", file)
	}
	return nil
}
