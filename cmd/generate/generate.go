package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// Config 配置参数
type Config struct {
	Dir     string
	PkgName string
	Name    string
	Comment string
	File    string
	Modules string
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

func (a *Command) hasModule(m string) bool {
	if v := a.cfg.Modules; v == "" || v == "all" {
		return true
	}

	for _, s := range strings.Split(a.cfg.Modules, ",") {
		if s == m {
			return true
		}
	}

	return false
}

func (a *Command) handleError(err error, desc string) {
	if err != nil {
		fmt.Printf("%s:%s", desc, err.Error())
	}
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

	if a.hasModule("schema") {
		err = genSchema(ctx, dir, item.StructName, item.Comment, item.toSchemaFields()...)
		if err != nil {
			return err
		}
	}

	if a.hasModule("entity") {
		err = genEntity(ctx, pkgName, dir, item.StructName, item.Comment, item.toEntityFields()...)
		a.handleError(err, "生成entity")
	}

	if a.hasModule("model") {
		err = genModelImpl(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成mode impl")

		err = genModel(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成model")

		err = insertModelInject(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成model inject")
	}

	if a.hasModule("bll") {
		err = genBllImpl(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成bll impl")

		err = genBll(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成bll")

		err = insertBllInject(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成bll inject")
	}

	if a.hasModule("ctl") {
		err = genCtl(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成ctl")

		err = insertCtlInject(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成ctl inject")
	}

	if a.hasModule("api") {
		err = insertAPI(ctx, pkgName, dir, item.StructName, item.Comment)
		a.handleError(err, "生成api")
	}

	return nil
}
