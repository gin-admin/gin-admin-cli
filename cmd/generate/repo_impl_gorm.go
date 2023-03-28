package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v6/util"
)

func getModelImplGormFileName(appName, dir, name string) string {
	name = util.ToLowerUnderlinedNamer(name)
	fullname := fmt.Sprintf("%s/internal/%s/dao/repo/%s.repo.go", dir, appName, name)
	return fullname
}

func genRepoImplGorm(ctx context.Context, obj *genObject) error {
	data := map[string]interface{}{
		"PkgName":       obj.pkgName,
		"AppName":       obj.appName,
		"Name":          obj.name,
		"PluralName":    util.ToPlural(obj.name),
		"Comment":       obj.comment,
		"UnderLineName": util.ToLowerUnderlinedNamer(obj.name),
		"IncludeStatus": !obj.excludeStatus,
	}

	buf, err := execParseTpl(daoGromRepoTpl, data)
	if err != nil {
		return err
	}

	fullname := getModelImplGormFileName(obj.appName, obj.dir, obj.name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const daoGromRepoTpl = `
package repo

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"{{.PkgName}}/internal/{{.AppName}}/dao/util"
	"{{.PkgName}}/internal/{{.AppName}}/schema"
	"{{.PkgName}}/pkg/errors"
)

// Injection wire
var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}Repo), "*"))

func Get{{.Name}}DB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(schema.{{.Name}}))
}

type {{.Name}}Repo struct {
	DB *gorm.DB
}

func (a *{{.Name}}Repo) Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error) {
	var opt schema.{{.Name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := Get{{.Name}}DB(ctx, a.DB)

	// TODO: Your where condition code here...

	if len(opt.SelectFields) > 0 {
		db = db.Select(opt.SelectFields)
	}

	if len(opt.OrderFields) > 0 {
		db = db.Order(util.ParseOrder(opt.OrderFields))
	}

	var list schema.{{.PluralName}}
	pr, err := util.WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.{{.Name}}QueryResult{
		PageResult: pr,
		Data:       list,
	}

	return qr, nil
}

func (a *{{.Name}}Repo) Get(ctx context.Context, id string) (*schema.{{.Name}}, error) {
	item := new(schema.{{.Name}})
	ok, err := util.FindOne(ctx, Get{{.Name}}DB(ctx, a.DB).Where("id=?", id), item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item, nil
}

func (a *{{.Name}}Repo) Create(ctx context.Context, item *schema.{{.Name}}) error {
	result := Get{{.Name}}DB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

func (a *{{.Name}}Repo) Update(ctx context.Context, item *schema.{{.Name}}) error {
	result := Get{{.Name}}DB(ctx, a.DB).Where("id=?", item.ID).Select("*").Updates(item)
	return errors.WithStack(result.Error)
}

func (a *{{.Name}}Repo) Delete(ctx context.Context, id string) error {
	result := Get{{.Name}}DB(ctx, a.DB).Where("id=?", id).Delete(new(schema.{{.Name}}))
	return errors.WithStack(result.Error)
}

{{if .IncludeStatus}}
func (a *{{.Name}}Repo) UpdateStatus(ctx context.Context, id string, status int) error {
	result := Get{{.Name}}DB(ctx, a.DB).Where("id=?", id).Update("status", status)
	return errors.WithStack(result.Error)
}
{{end}}

`
