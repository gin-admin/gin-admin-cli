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
}

func NewGenerate(cfg *GenerateConfig) *Generate {
	return &Generate{
		logger:           zap.S().Named("[Gen]"),
		cfg:              cfg,
		fs:               tfs.Ins,
		rootImportPath:   parser.GetRootImportPath(cfg.Dir),
		moduleImportPath: parser.GetModuleImportPath(cfg.Dir, cfg.ModulePath, cfg.ModuleName),
	}
}

type Generate struct {
	logger           *zap.SugaredLogger
	cfg              *GenerateConfig
	fs               tfs.FS
	rootImportPath   string
	moduleImportPath string
}

// Run generate command
func (a *Generate) Run(ctx context.Context, configFile string) error {
	switch filepath.Ext(configFile) {
	case ".json":
		var data []*schema.S
		if err := utils.ParseJSONFile(configFile, &data); err != nil {
			return err
		}
		return a.run(ctx, data)
	case ".yaml", "yml":
		var data []*schema.S
		if err := utils.ParseYAMLFile(configFile, &data); err != nil {
			return err
		}
		return a.run(ctx, data)
	default:
		return fmt.Errorf("unsupported config file type: %s", configFile)
	}
}

func (a *Generate) RunWithStruct(ctx context.Context, structName, comment, output string) error {
	return a.run(ctx, []*schema.S{
		{Name: structName, Comment: comment, Outputs: strings.Split(output, ",")},
	})
}

func (a *Generate) run(ctx context.Context, data []*schema.S) error {
	for _, d := range data {
		err := a.generate(ctx, d)
		if err != nil {
			return err
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

func (a *Generate) getGoTplFile(pkgName string) string {
	pkgName = fmt.Sprintf("%s.go.tpl", pkgName)
	if a.cfg.TplType != "" {
		return filepath.Join(a.cfg.TplType, pkgName)
	}
	return pkgName
}

func (a Generate) getAbsPath(file string) (string, error) {
	modPath := a.cfg.ModulePath
	file = filepath.Join(a.cfg.Dir, modPath, file)
	fullpath, err := filepath.Abs(file)
	if err != nil {
		a.logger.Errorf("Failed to get abs path, err: %s, #file %s", err, file)
		return "", err
	}
	return fullpath, nil
}

func (a *Generate) write(ctx context.Context, moduleName, structName, tpl string, data []byte, checkExists bool) error {
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

	if err := utils.ExecGoImports(file); err != nil {
		a.logger.Errorf("Failed to exec go imports, err: %s, #file %s", err, file)
		return nil
	}
	return nil
}

func (a *Generate) generate(ctx context.Context, dataItem *schema.S) error {
	dataItem = dataItem.Format()
	dataItem.RootImportPath = a.rootImportPath
	dataItem.ModuleName = a.cfg.ModuleName
	dataItem.ModuleImportPath = a.moduleImportPath

	genPackages := parser.StructPackages
	if len(dataItem.Outputs) > 0 {
		genPackages = dataItem.Outputs
	}

	for _, pkgName := range genPackages {
		tplName := a.getGoTplFile(pkgName)
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

func (a *Generate) execWireAndSwag(ctx context.Context) error {
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
