package generate

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getBllImplFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/bll/impl/internal/b_%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成bll实现文件
func genBllImpl(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName": pkgName,
		"Name":    name,
		"Comment": comment,
	}

	buf, err := execParseTpl(bllImplTpl, data)
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

const bllImplTpl = `
package internal

import (
	"context"

	"{{.PkgName}}/internal/app/errors"
	"{{.PkgName}}/internal/app/model"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/util"
)

// New{{.Name}} 创建{{.Comment}}
func New{{.Name}}(m{{.Name}} model.I{{.Name}}) *{{.Name}} {
	return &{{.Name}}{
		{{.Name}}Model: m{{.Name}},
	}
}

// {{.Name}} {{.Comment}}业务逻辑
type {{.Name}} struct {
	{{.Name}}Model model.I{{.Name}}
}

// Query 查询数据
func (a *{{.Name}}) Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error) {
	return a.{{.Name}}Model.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *{{.Name}}) Get(ctx context.Context, recordID string, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}, error) {
	item, err := a.{{.Name}}Model.Get(ctx, recordID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *{{.Name}}) getUpdate(ctx context.Context, recordID string) (*schema.{{.Name}}, error) {
	return a.Get(ctx, recordID)
}

// Create 创建数据
func (a *{{.Name}}) Create(ctx context.Context, item schema.{{.Name}}) (*schema.{{.Name}}, error) {
	item.RecordID = util.MustUUID()
	err := a.{{.Name}}Model.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.RecordID)
}

// Update 更新数据
func (a *{{.Name}}) Update(ctx context.Context, recordID string, item schema.{{.Name}}) (*schema.{{.Name}}, error) {
	oldItem, err := a.{{.Name}}Model.Get(ctx, recordID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	}

	err = a.{{.Name}}Model.Update(ctx, recordID, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, recordID)
}

// Delete 删除数据
func (a *{{.Name}}) Delete(ctx context.Context, recordID string) error {
	oldItem, err := a.{{.Name}}Model.Get(ctx, recordID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.{{.Name}}Model.Delete(ctx, recordID)
}

`
