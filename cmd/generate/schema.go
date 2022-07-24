package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

type schemaField struct {
	Name           string
	Comment        string
	Type           string
	IsRequired     bool
	BindingOptions string
}

func getSchemaFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/schema/%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genSchema(ctx context.Context, pkgName, dir string, item TplItem, excludeStatus, excludeCreate bool, fields ...schemaField) error {
	var tfields []schemaField

	tfields = append(tfields, schemaField{Name: "ID", Type: "uint64"})
	tfields = append(tfields, fields...)

	if !excludeStatus {
		tfields = append(tfields, schemaField{Name: "Status", Comment: "1:enable,2:disable", Type: "int"})
	}

	if !excludeCreate {
		tfields = append(tfields, schemaField{Name: "Creator", Type: "uint64"})
	}

	tfields = append(tfields, schemaField{Name: "CreatedAt", Type: "time.Time"})
	tfields = append(tfields, schemaField{Name: "UpdatedAt", Type: "time.Time"})

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
		"Name":       item.StructName,
		"PluralName": util.ToPlural(item.StructName),
		"Fields":     buf.String(),
		"Comment":    item.Comment,
		"ItemFields": item.Fields,
	})
	if err != nil {
		return err
	}

	fullname := getSchemaFileName(dir, item.StructName)
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
	PaginationParam` +
	"{{range .ItemFields}}" +
	"{{if .ConditionArray}}" +
	"\n{{fieldToPlural .StructFieldName}}  []{{.StructFieldType}}         `form:\"{{fieldToLowerUnderlinedName .StructFieldName}}\"` // {{.Comment}}" +
	"{{end}}" +
	"{{if .Condition}}" +
	"\n{{.StructFieldName}}  {{.StructFieldType}}        `form:\"{{fieldToLowerUnderlinedName .StructFieldName}}\"`// {{.Comment}}" +
	"{{end}}" +
	"{{if .ConditionLike}}" +
	"\n{{.StructFieldName}}Like  {{.StructFieldType}} `form:\"{{fieldToLowerUnderlinedName .StructFieldName}}_like\"`// {{.Comment}}" +
	"{{end}}" +
	"{{end}}" + `
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
