package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

type schemaField struct {
	Name           string // 字段名
	Comment        string // 字段注释
	Type           string // 字段类型
	IsRequired     bool   // 是否必选项
	BindingOptions string // binding配置项(不包含required，required由IsRequired控制)
}

func getSchemaFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/schema/%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genSchema(ctx context.Context, pkgName, dir, name, comment string, excludeStatus, excludeCreate bool, fields ...schemaField) error {
	var tfields []schemaField

	tfields = append(tfields, schemaField{Name: "ID", Comment: "唯一标识", Type: "uint64"})
	tfields = append(tfields, fields...)

	if !excludeStatus {
		tfields = append(tfields, schemaField{Name: "Status", Comment: "状态(1:启用 2:禁用)", Type: "int"})
	}

	if !excludeCreate {
		tfields = append(tfields, schemaField{Name: "Creator", Comment: "创建者", Type: "uint64"})
	}

	tfields = append(tfields, schemaField{Name: "CreatedAt", Comment: "创建时间", Type: "time.Time"})
	tfields = append(tfields, schemaField{Name: "UpdatedAt", Comment: "更新时间", Type: "time.Time"})

	buf := new(bytes.Buffer)
	for _, field := range tfields {
		buf.WriteString(fmt.Sprintf("%s \t %s \t", field.Name, field.Type))
		buf.WriteByte('`')
		if field.Name == "ID" {
			buf.WriteString(fmt.Sprintf(`json:"%s,string"`, util.ToLowerUnderlinedNamer(field.Name)))
		} else {
			buf.WriteString(fmt.Sprintf(`json:"%s"`, util.ToLowerUnderlinedNamer(field.Name)))
		}

		bindingOpts := ""
		if field.IsRequired {
			bindingOpts = "required"
		}

		if v := field.BindingOptions; v != "" {
			if bindingOpts != "" {
				bindingOpts += ","
			}
			bindingOpts = bindingOpts + v
		}
		if bindingOpts != "" {
			buf.WriteByte(' ')
			buf.WriteString(fmt.Sprintf(`binding:"%s"`, bindingOpts))
		}

		buf.WriteByte('`')

		if field.Comment != "" {
			buf.WriteString(fmt.Sprintf("// %s", field.Comment))
		}
		buf.WriteString(delimiter)
	}

	tbuf, err := execParseTpl(schemaTpl, map[string]interface{}{
		"PkgName":    pkgName,
		"Name":       name,
		"PluralName": util.ToPlural(name),
		"Fields":     buf.String(),
		"Comment":    comment,
	})
	if err != nil {
		return err
	}

	fullname := getSchemaFileName(dir, name)
	err = createFile(ctx, fullname, tbuf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const schemaTpl = `
package schema

import "time"

// {{.Comment}}
type {{.Name}} struct {
	{{.Fields}}
}

// Query parameters for db
type {{.Name}}QueryParam struct {
	PaginationParam
}

// Query options for db (order or select fields)
type {{.Name}}QueryOptions struct {
	OrderFields []*OrderField
	SelectFields []string
}

// Query result from db
type {{.Name}}QueryResult struct {
	Data       {{.PluralName}}
	PageResult *PaginationResult
}

// {{.Comment}} Object List
type {{.PluralName}} []*{{.Name}}

`
