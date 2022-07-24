package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getWebPageFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/service/%s.srv.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genWebPage(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName": pkgName,
		"Name":    name,
		"Comment": comment,
	}

	buf, err := execParseTpl(webPageTpl, data)
	if err != nil {
		return err
	}

	fullname := getWebPageFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const webPageTpl = `
package service

import (
	"context"

	"{{.PkgName}}/internal/app/schema"
)

// I{{.Name}} {{.Comment}}业务逻辑接口
type I{{.Name}} interface {
	// 查询数据
	Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, id string, opts ...schema.{{.Name}}GetOptions) (*schema.{{.Name}}, error)
	// 创建数据
	Create(ctx context.Context, item schema.{{.Name}}) (*schema.IDResult, error)
	// 更新数据
	Update(ctx context.Context, id string, item schema.{{.Name}}) error
	// 删除数据
	Delete(ctx context.Context, id string) error
}

`
