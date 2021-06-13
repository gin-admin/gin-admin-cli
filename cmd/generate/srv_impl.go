package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v4/util"
)

func getBllImplFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/service/%s.srv.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成bll实现文件
func genBllImpl(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName": pkgName,
		"Name":    name,
		"Comment": comment,
	}

	buf, err := execParseTpl(srvImplTpl, data)
	if err != nil {
		return err
	}

	fullname := getBllImplFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const srvImplTpl = `
package service

import (
	"context"

	"{{.PkgName}}/internal/app/model/gormx/entity"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/errors"
	"{{.PkgName}}/pkg/util/uuid"
	"github.com/google/wire"
)

// var _ srv.I{{.Name}} = (*{{.Name}})(nil)

// {{.Name}}Set 注入{{.Name}}
var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}), "*"))

// {{.Name}} {{.Comment}}
type {{.Name}} struct {
	{{.Name}}Model *entity.{{.Name}}
}

// Query 查询数据
func (a *{{.Name}}) Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error) {
	return a.{{.Name}}Model.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *{{.Name}}) Get(ctx context.Context, id string, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}, error) {
	item, err := a.{{.Name}}Model.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

// Create 创建数据
func (a *{{.Name}}) Create(ctx context.Context, item schema.{{.Name}}) (*schema.IDResult, error) {
	// TODO: check?
	item.ID = uuid.MustString()
	err := a.{{.Name}}Model.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

// Update 更新数据
func (a *{{.Name}}) Update(ctx context.Context, id string, item schema.{{.Name}}) error {
	oldItem, err := a.{{.Name}}Model.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	// TODO: check?
	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	return a.{{.Name}}Model.Update(ctx, id, item)
}

// Delete 删除数据
func (a *{{.Name}}) Delete(ctx context.Context, id string) error {
	oldItem, err := a.{{.Name}}Model.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.{{.Name}}Model.Delete(ctx, id)
}

`
