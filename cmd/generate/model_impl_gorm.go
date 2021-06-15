package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v4/util"
)

func getModelImplGormFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/gormx/repo/%s.repo.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成model实现文件
func genModelImplGorm(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName":    pkgName,
		"Name":       name,
		"PluralName": util.ToPlural(name),
		"Comment":    comment,
	}

	buf, err := execParseTpl(modelImplGormTpl, data)
	if err != nil {
		return err
	}

	fullname := getModelImplGormFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const modelImplGormTpl = `
package repo

import (
	"context"

	"{{.PkgName}}/internal/app/model/gormx/entity"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/errors"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

// var _ model.I{{.Name}} = (*{{.Name}})(nil)

// {{.Name}}Set 注入{{.Name}}
var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}), "*"))

// {{.Name}} {{.Comment}}存储
type {{.Name}} struct {
	DB *gorm.DB
}

func (a *{{.Name}}) getQueryOption(opts ...schema.{{.Name}}QueryOptions) schema.{{.Name}}QueryOptions {
	var opt schema.{{.Name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *{{.Name}}) Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.Get{{.Name}}DB(ctx, a.DB)
	// TODO: 查询条件

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.{{.PluralName}}
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.{{.Name}}QueryResult{
		PageResult: pr,
		Data:       list.ToSchema{{.PluralName}}(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *{{.Name}}) Get(ctx context.Context, id string, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}, error) {
	db := entity.Get{{.Name}}DB(ctx, a.DB).Where("id=?", id)
	var item entity.{{.Name}}
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchema{{.Name}}(), nil
}

// Create 创建数据
func (a *{{.Name}}) Create(ctx context.Context, item schema.{{.Name}}) error {
	eitem := entity.Schema{{.Name}}(item).To{{.Name}}()
	result := entity.Get{{.Name}}DB(ctx, a.DB).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *{{.Name}}) Update(ctx context.Context, id string, item schema.{{.Name}}) error {
	eitem := entity.Schema{{.Name}}(item).To{{.Name}}()
	result := entity.Get{{.Name}}DB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *{{.Name}}) Delete(ctx context.Context, id string) error {
	result := entity.Get{{.Name}}DB(ctx, a.DB).Where("id=?", id).Delete(entity.{{.Name}}{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

`
