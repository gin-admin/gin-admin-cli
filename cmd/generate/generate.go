package generate

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Dir           string
	PkgName       string
	Name          string
	Comment       string
	File          string
	Modules       string
	ExcludeStatus bool
	ExcludeCreate bool
}

func Exec(cfg Config) error {
	cmd := &Command{cfg: &cfg}
	return cmd.Exec()
}

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
		fmt.Printf("%s: %s", desc, err.Error())
	}
}

type genObject struct {
	pkgName       string
	appName       string
	dir           string
	name          string
	comment       string
	excludeStatus bool
	excludeCreate bool
	fields        []schemaField
}

func (a *Command) Exec() error {
	var item TplItem

	if a.cfg.File != "" {
		b, err := readFile(a.cfg.File)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(b.Bytes(), &item)
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

	ctx := context.Background()
	genObj := &genObject{
		pkgName:       a.cfg.PkgName,
		appName:       filepath.Base(a.cfg.PkgName),
		dir:           dir,
		name:          item.StructName,
		comment:       item.Comment,
		excludeStatus: a.cfg.ExcludeStatus,
		excludeCreate: a.cfg.ExcludeCreate,
		fields:        item.toSchemaFields(),
	}

	if a.hasModule("schema") {
		err = genSchema(ctx, genObj)
		if err != nil && err != ErrFileExists {
			return err
		}
	}

	if a.hasModule("dao") {
		err = genRepoImplGorm(ctx, genObj)
		a.handleError(err, "Generate gorm model")

		err = insertRepoInjectGorm(ctx, genObj)
		a.handleError(err, "Insert gorm model inject")
	}

	if a.hasModule("service") {
		err = genServiceImpl(ctx, genObj)
		a.handleError(err, "Generate bll impl")

		err = insertServiceInject(ctx, genObj)
		a.handleError(err, "Insert bll inject")
	}

	if a.hasModule("api") {
		err = genAPI(ctx, genObj)
		a.handleError(err, "Generate api")

		err = insertAPIInject(ctx, genObj)
		a.handleError(err, "Insert api inject")
	}

	if a.hasModule("router") {
		err = insertRouterAPI(ctx, genObj)
		a.handleError(err, "Insert router api")

		err = insertRouterInject(ctx, genObj)
		a.handleError(err, "Insert router inject")
	}

	return nil
}
