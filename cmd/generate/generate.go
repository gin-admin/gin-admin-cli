package generate

import (
	"context"
	"encoding/json"
	"path/filepath"
)

// Config 配置参数
type Config struct {
	Dir     string
	PkgName string
	Name    string
	Comment string
	File    string
}

// Exec 执行生成模块命令
func Exec(cfg Config) error {
	cmd := &Command{cfg: &cfg}
	return cmd.Exec()
}

// Command 生成命令
type Command struct {
	cfg *Config
}

// Exec 执行命令
func (a *Command) Exec() error {
	var item TplItem

	if a.cfg.File != "" {
		b, err := readFile(a.cfg.File)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b.Bytes(), &item)
		if err != nil {
			return err
		}
	} else {
		item.StructName = a.cfg.Name
		item.Comment = a.cfg.Comment
	}

	dir, err := filepath.Abs(a.cfg.Dir)
	if err != nil {
		return err
	}

	pkgName := a.cfg.PkgName

	ctx := context.Background()
	err = genSchema(ctx, dir, item.StructName, item.Comment, item.toSchemaFields()...)
	if err != nil {
		return err
	}

	err = genEntity(ctx, pkgName, dir, item.StructName, item.Comment, item.toEntityFields()...)
	if err != nil {
		return err
	}

	err = genModelImpl(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = genModel(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = insertModelInject(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = genBllImpl(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = genBll(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = insertBllInject(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = genCtl(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = insertCtlInject(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	err = insertAPI(ctx, pkgName, dir, item.StructName, item.Comment)
	if err != nil {
		return err
	}

	return nil
}
