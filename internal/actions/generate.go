package actions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-admin/gin-admin-cli/v10/internal/parser"
	"github.com/gin-admin/gin-admin-cli/v10/internal/schema"
	"github.com/gin-admin/gin-admin-cli/v10/internal/tfs"
	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
	"go.uber.org/zap"
)

type GenerateConfig struct {
	Dir         string
	TplType     string
	Module      string
	ModulePath  string
	WirePath    string
	SwaggerPath string
	FEDir       string
}

func Generate(cfg GenerateConfig) *GenerateAction {
	return &GenerateAction{
		logger:           zap.S().Named("[GEN]"),
		cfg:              &cfg,
		fs:               tfs.Ins,
		rootImportPath:   parser.GetRootImportPath(cfg.Dir),
		moduleImportPath: parser.GetModuleImportPath(cfg.Dir, cfg.ModulePath, cfg.Module),
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
	var parseFile = func(name string) ([]*schema.S, error) {
		var data []*schema.S
		switch filepath.Ext(name) {
		case ".json":
			if err := utils.ParseJSONFile(name, &data); err != nil {
				return nil, err
			}
		case ".yaml", ".yml":
			if err := utils.ParseYAMLFile(name, &data); err != nil {
				return nil, err
			}
		default:
			a.logger.Warnf("Ignore file %s, only support json/yaml/yml", name)
		}
		if len(data) == 0 {
			return nil, nil
		}
		return data, nil
	}

	if utils.IsDir(cfgName) {
		var data []*schema.S
		err := filepath.WalkDir(cfgName, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			items, err := parseFile(path)
			if err != nil {
				return err
			}
			data = append(data, items...)
			return nil
		})
		if err != nil {
			return err
		}
		return a.run(ctx, data)
	}

	data, err := parseFile(cfgName)
	if err != nil {
		return err
	} else if len(data) == 0 {
		a.logger.Warnf("No data found in file %s", cfgName)
		return nil
	}
	return a.run(ctx, data)
}

func (a *GenerateAction) RunWithStruct(ctx context.Context, s *schema.S) error {
	return a.run(ctx, []*schema.S{s})
}

func (a *GenerateAction) run(ctx context.Context, data []*schema.S) error {
	moduleMap := make(map[string]bool)

	for _, d := range data {
		if d.Module == "" && a.cfg.Module == "" {
			return fmt.Errorf("Struct %s module is empty", d.Name)
		}

		if d.Module == "" {
			d.Module = a.cfg.Module
		}
		if !moduleMap[d.Module] {
			moduleMap[d.Module] = true
		}
		if err := a.generate(ctx, d); err != nil {
			return err
		}
		if d.GenerateFE {
			if err := a.generateFE(ctx, d); err != nil {
				return err
			}
		}
	}

	for module := range moduleMap {
		modsTplData, err := parser.ModifyModsFile(ctx, parser.BasicArgs{
			Dir:        a.cfg.Dir,
			ModuleName: module,
			ModulePath: a.cfg.ModulePath,
			Flag:       parser.AstFlagGen,
		})
		if err != nil {
			a.logger.Errorf("Failed to modify mods file, err: %s", err)
			return err
		}

		err = a.write(ctx, module, "", parser.FileForMods, modsTplData, false)
		if err != nil {
			return err
		}
	}

	return a.execWireAndSwag(ctx)
}

func (a *GenerateAction) getGoTplFile(tplName, tplType string) string {
	tplName = fmt.Sprintf("%s.go.tpl", tplName)
	if tplType == "" && a.cfg.TplType != "" {
		tplType = a.cfg.TplType
	}

	if tplType != "" {
		p := filepath.Join(tplType, tplName)
		if ok, _ := utils.ExistsFile(p); ok {
			return p
		}
		return filepath.Join("default", tplName)
	}
	return tplName
}

func (a GenerateAction) getAbsPath(file string) (string, error) {
	modPath := a.cfg.ModulePath
	file = filepath.Join(a.cfg.Dir, modPath, file)
	fullPath, err := filepath.Abs(file)
	if err != nil {
		a.logger.Errorf("Failed to get abs path, err: %s, #file %s", err, file)
		return "", err
	}
	return fullPath, nil
}

func (a *GenerateAction) write(_ context.Context, moduleName, structName, tpl string, data []byte, checkExists bool) error {
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

		var rewrite bool
		switch pkgName {
		case "schema":
			if dataItem.Rewrite != nil && dataItem.Rewrite.Schema {
				rewrite = true
			}
		case "dal":
			if dataItem.Rewrite != nil && dataItem.Rewrite.Dal {
				rewrite = true
			}
		case "biz":
			if dataItem.Rewrite != nil && dataItem.Rewrite.Biz {
				rewrite = true
			}
		case "api":
			if dataItem.Rewrite != nil && dataItem.Rewrite.Api {
				rewrite = true
			}
		}

		err = a.write(ctx, dataItem.Module, dataItem.Name, parser.StructPackageTplPaths[pkgName], tplData, !rewrite)
		if err != nil {
			return err
		}
	}

	basicArgs := parser.BasicArgs{
		Dir:              a.cfg.Dir,
		ModuleName:       dataItem.Module,
		ModulePath:       a.cfg.ModulePath,
		StructName:       dataItem.Name,
		GenPackages:      genPackages,
		Flag:             parser.AstFlagGen,
		FillRouterPrefix: dataItem.FillRouterPrefix,
	}
	moduleMainTplData, err := parser.ModifyModuleMainFile(ctx, basicArgs)
	if err != nil {
		a.logger.Errorf("Failed to modify module main file, err: %s, #struct %s", err, dataItem.Name)
		return err
	}

	err = a.write(ctx, dataItem.Module, dataItem.Name, parser.FileForModuleMain, moduleMainTplData, false)
	if err != nil {
		return err
	}

	moduleWireTplData, err := parser.ModifyModuleWireFile(ctx, basicArgs)
	if err != nil {
		a.logger.Errorf("Failed to modify module wire file, err: %s, #struct %s", err, dataItem.Name)
		return err
	}

	err = a.write(ctx, dataItem.Module, dataItem.Name, parser.FileForModuleWire, moduleWireTplData, false)
	if err != nil {
		return err
	}

	return nil
}

func (a *GenerateAction) execWireAndSwag(_ context.Context) error {
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

func (a *GenerateAction) generateFE(_ context.Context, dataItem *schema.S) error {
	for tpl, file := range dataItem.FEMapping {
		tplPath := filepath.Join(dataItem.FETpl, tpl)
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
