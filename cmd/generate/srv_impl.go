package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v6/util"
)

func getBllImplFileName(appName, dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/%s/service/%s.srv.go", dir, appName, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genServiceImpl(ctx context.Context, obj *genObject) error {
	data := map[string]interface{}{
		"PkgName":       obj.pkgName,
		"AppName":       obj.appName,
		"Name":          obj.name,
		"Comment":       obj.comment,
		"IncludeStatus": !obj.excludeStatus,
		"IncludeCreate": !obj.excludeCreate,
	}

	buf, err := execParseTpl(serviceImplTpl, data)
	if err != nil {
		return err
	}

	fullname := getBllImplFileName(obj.appName, obj.dir, obj.name)
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
	"fmt"
	"context"

	"github.com/google/wire"

	"{{.PkgName}}/internal/{{.AppName}}/dao"
	"{{.PkgName}}/internal/{{.AppName}}/schema"
	"{{.PkgName}}/internal/{{.AppName}}/module/consts"
	{{if .IncludeCreate}}"{{.PkgName}}/internal/{{.AppName}}/module/contextx"{{end}}
	"{{.PkgName}}/pkg/errors"
	"{{.PkgName}}/pkg/util/xid"
)

var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}Srv), "*"))

type {{.Name}}Srv struct {
	TransRepo              		*dao.TransRepo
	{{.Name}}Repo               *dao.{{.Name}}Repo
}

func (a *{{.Name}}Srv) Query(ctx context.Context, params schema.{{.Name}}QueryParam) (*schema.{{.Name}}QueryResult, error) {
	params.Pagination = true
	result, err := a.{{.Name}}Repo.Query(ctx, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *{{.Name}}Srv) Get(ctx context.Context, id string) (*schema.{{.Name}}, error) {
	item, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.NotFound(consts.ErrNotFoundID, "{{.Name}} not found")
	}

	return item, nil
}

func (a *{{.Name}}Srv) Create(ctx context.Context, item schema.{{.Name}}) (*schema.{{.Name}}, error) {
	item.ID = xid.NewID()

	{{if .IncludeCreate}}
	item.CreatedBy = contextx.FromUserID(ctx)
	{{end}}

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.{{.Name}}Repo.Create(ctx, &item)
	})
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (a *{{.Name}}Srv) Update(ctx context.Context, id string, item schema.{{.Name}}) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.NotFound(consts.ErrNotFoundID, "{{.Name}} not found")
	}

	item.ID = id
	item.CreatedAt = oldItem.CreatedAt
	{{if .IncludeCreate}}
	item.CreatedBy = oldItem.CreatedBy
	{{end}}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.{{.Name}}Repo.Update(ctx, &item)
	})
}

func (a *{{.Name}}Srv) Delete(ctx context.Context, id string) error {
	oldItem, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.NotFound(consts.ErrNotFoundID, "{{.Name}} not found")
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.{{.Name}}Repo.Delete(ctx, id)
	})
}

{{if .IncludeStatus}}
func (a *{{.Name}}Srv) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.{{.Name}}Repo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.NotFound(consts.ErrNotFoundID, "{{.Name}} not found")
	} else if oldItem.Status == status {
		return nil
	}

	return a.{{.Name}}Repo.UpdateStatus(ctx, id, status)
}
{{end}}

`
