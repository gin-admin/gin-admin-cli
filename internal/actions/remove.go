package actions

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/parser"
	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
	"go.uber.org/zap"
)

type RemoveConfig struct {
	Dir         string
	ModuleName  string
	ModulePath  string
	WirePath    string
	SwaggerPath string
}

func NewRemove(cfg *RemoveConfig) *Remove {
	return &Remove{
		logger: zap.S().Named("[Rem]"),
		cfg:    cfg,
	}
}

type Remove struct {
	logger *zap.SugaredLogger
	cfg    *RemoveConfig
}

func (a *Remove) Run(ctx context.Context, structs string) error {
	for _, name := range strings.Split(structs, ",") {
		for _, pkgName := range parser.StructPackages {
			err := a.modify(ctx, a.cfg.ModuleName, name, parser.StructPackageTplPaths[pkgName], nil, true)
			if err != nil {
				return err
			}
		}

		basicArgs := parser.BasicArgs{
			Dir:        a.cfg.Dir,
			ModuleName: a.cfg.ModuleName,
			ModulePath: a.cfg.ModulePath,
			StructName: name,
			Flag:       parser.AstFlagRem,
		}
		moduleMainTplData, err := parser.ModifyModuleMainFile(ctx, basicArgs)
		if err != nil {
			a.logger.Errorf("Failed to modify module main file, err: %s, #struct %s", err, name)
			return err
		}

		err = a.modify(ctx, a.cfg.ModuleName, name, parser.FileForModuleMain, moduleMainTplData, false)
		if err != nil {
			return err
		}

		moduleWireTplData, err := parser.ModifyModuleWireFile(ctx, basicArgs)
		if err != nil {
			a.logger.Errorf("Failed to modify module wire file, err: %s, #struct %s", err, name)
			return err
		}

		err = a.modify(ctx, a.cfg.ModuleName, name, parser.FileForModuleWire, moduleWireTplData, false)
		if err != nil {
			return err
		}
	}

	return a.execWireAndSwag(ctx)
}

func (a Remove) getAbsPath(file string) (string, error) {
	modPath := a.cfg.ModulePath
	file = filepath.Join(a.cfg.Dir, modPath, file)
	fullpath, err := filepath.Abs(file)
	if err != nil {
		a.logger.Errorf("Failed to get abs path, err: %s, #file %s", err, file)
		return "", err
	}
	return fullpath, nil
}

func (a *Remove) modify(ctx context.Context, moduleName, structName, tpl string, data []byte, deleted bool) error {
	file, err := parser.ParseFilePathFromTpl(moduleName, structName, tpl)
	if err != nil {
		a.logger.Errorf("Failed to parse file path from tpl, err: %s, #struct %s, #tpl %s", err, structName, tpl)
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

	if exists {
		if err := os.Remove(file); err != nil {
			a.logger.Errorf("Failed to remove file, err: %s, #file %s", err, file)
			return err
		}
	}

	if deleted {
		a.logger.Infof("Delete file: %s", file)
		return nil
	}

	if !exists {
		return nil
	}

	a.logger.Infof("Write file: %s", file)
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

func (a *Remove) execWireAndSwag(ctx context.Context) error {
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
