package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getBllImplFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/service/%s.srv.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genServiceImpl(ctx context.Context, pkgName, dir, name, comment string, excludeStatus, excludeCreate bool) error {
	data := map[string]interface{}{
		"PkgName":       pkgName,
		"Name":          name,
		"Comment":       comment,
		"IncludeStatus": !excludeStatus,
		"IncludeCreate": !excludeCreate,
	}

	buf, err := execParseTpl(serviceImplTpl, data)
	if err != nil {
		return err
	}

	fullname := getBllImplFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const serviceImplTpl = `
package service

import (
	"context"

	"github.com/google/wire"

	"{{.PkgName}}/internal/app/dao"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/errors"
	"{{.PkgName}}/pkg/util/snowflake"
)

var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}Srv), "*"))

type {{.Name}}Srv struct {
	TransRepo              		*dao.TransRepo
	{{.Name}}Repo               *dao.{{.Name}}Repo
}

func (a *{{.Name}}Srv) Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error) {
	result, err := a.{{.Name}}Repo.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *{{.Name}}Srv) Get(ctx context.Context, id uint64, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}, error) {
	item, err := a.{{.Name}}Repo.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *{{.Name}}Srv) Create(ctx context.Context, item schema.{{.Name}}) (*schema.IDResult, error) {
	item.ID = snowflake.MustID()

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.{{.Name}}Repo.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *{{.Name}}Srv) Update(ctx context.Context, id uint64, item schema.{{.Name}}) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	item.ID = oldItem.ID
	item.CreatedAt = oldItem.CreatedAt
	{{if .IncludeCreate}}
	item.Creator = oldItem.Creator
	{{end}}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.{{.Name}}Repo.Update(ctx, id, item)
	})
}

func (a *{{.Name}}Srv) Delete(ctx context.Context, id uint64) error {
	oldItem, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.{{.Name}}Repo.Delete(ctx, id)
	})
}

{{if .IncludeStatus}}
func (a *{{.Name}}Srv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	oldItem, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Status == status {
		return nil
	}

	return a.{{.Name}}Repo.UpdateStatus(ctx, id, status)
}
{{end}}

`
