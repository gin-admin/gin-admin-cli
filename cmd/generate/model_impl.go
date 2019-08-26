package generate

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getModelImplFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/gorm/internal/model/m_%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成model实现文件
func genModelImpl(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName": pkgName,
		"Name":    name,
		"Comment": comment,
	}

	buf, err := execParseTpl(modelImplTpl, data)
	if err != nil {
		return err
	}

	fullname := getModelImplFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const modelImplTpl = `
package model

import (
	"context"

	"{{.PkgName}}/internal/app/errors"
	"{{.PkgName}}/internal/app/model/impl/gorm/internal/entity"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/gormplus"
)

// New{{.Name}} 创建{{.Comment}}存储实例
func New{{.Name}}(db *gormplus.DB) *{{.Name}} {
	return &{{.Name}}{db}
}

// {{.Name}} {{.Comment}}存储
type {{.Name}} struct {
	db *gormplus.DB
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
	db := entity.Get{{.Name}}DB(ctx, a.db).DB
	
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.{{.Name}}s
	pr, err := WrapPageQuery(db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.{{.Name}}QueryResult{
		PageResult: pr,
		Data:       list.ToSchema{{.Name}}s(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *{{.Name}}) Get(ctx context.Context, recordID string, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}, error) {
	db := entity.Get{{.Name}}DB(ctx, a.db).Where("record_id=?", recordID)
	var item entity.{{.Name}}
	ok, err := a.db.FindOne(db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchema{{.Name}}(), nil
}

// Create 创建数据
func (a *{{.Name}}) Create(ctx context.Context, item schema.{{.Name}}) error {
	{{.Name}} := entity.Schema{{.Name}}(item).To{{.Name}}()
	result := entity.Get{{.Name}}DB(ctx, a.db).Create({{.Name}})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *{{.Name}}) Update(ctx context.Context, recordID string, item schema.{{.Name}}) error {
	{{.Name}} := entity.Schema{{.Name}}(item).To{{.Name}}()
	result := entity.Get{{.Name}}DB(ctx, a.db).Where("record_id=?", recordID).Omit("record_id", "creator").Updates({{.Name}})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *{{.Name}}) Delete(ctx context.Context, recordID string) error {
	result := entity.Get{{.Name}}DB(ctx, a.db).Where("record_id=?", recordID).Delete(entity.{{.Name}}{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

`
