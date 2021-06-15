package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v4/util"
)

func getModelFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/%s.repo.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成model文件
func genModel(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName": pkgName,
		"Name":    name,
		"Comment": comment,
	}

	buf, err := execParseTpl(modelTpl, data)
	if err != nil {
		return err
	}

	fullname := getModelFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const modelTpl = `
package repo

import (
	"context"

	"{{.PkgName}}/internal/app/schema"
)

// I{{.Name}} {{.Comment}}存储接口
type I{{.Name}} interface {
	// 查询数据
	Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, id string, opts ...schema.{{.Name}}GetOptions) (*schema.{{.Name}}, error)
	// 创建数据
	Create(ctx context.Context, item schema.{{.Name}}) error
	// 更新数据
	Update(ctx context.Context, id string, item schema.{{.Name}}) error
	// 删除数据
	Delete(ctx context.Context, id string) error
}

`
