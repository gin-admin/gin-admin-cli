package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v6/util"
)

type schemaField struct {
	Name           string
	Comment        string
	Type           string
	IsRequired     bool
	BindingOptions string
	GormOptions    string
}

func getSchemaFileName(appName, dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/%s/schema/%s.go", dir, appName, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genSchema(ctx context.Context, obj *genObject) error {
	var tfields []schemaField

	tfields = append(tfields, schemaField{Name: "ID", Type: "string", GormOptions: "size:20;primarykey;"})
	tfields = append(tfields, obj.fields...)

	if !obj.excludeStatus {
		tfields = append(tfields, schemaField{Name: "Status", Comment: "1:enable,2:disable", Type: "int"})
	}

	if !obj.excludeCreate {
		tfields = append(tfields, schemaField{Name: "CreatedBy", Type: "string", GormOptions: "size:20;"})
	}

	tfields = append(tfields, schemaField{Name: "CreatedAt", Type: "time.Time"})
	tfields = append(tfields, schemaField{Name: "UpdatedAt", Type: "time.Time"})

	buf := new(bytes.Buffer)
	for _, field := range tfields {
		buf.WriteString(fmt.Sprintf("%s \t %s \t", field.Name, field.Type))
		buf.WriteByte('`')

		if field.GormOptions != "" {
			buf.WriteString(fmt.Sprintf(`gorm:"%s"`, field.GormOptions))
			buf.WriteByte(' ')
		}

		buf.WriteString(fmt.Sprintf(`json:"%s"`, util.ToLowerUnderlinedNamer(field.Name)))

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
		"PkgName":    obj.pkgName,
		"Name":       obj.name,
		"PluralName": util.ToPlural(obj.name),
		"Fields":     buf.String(),
		"Comment":    obj.comment,
	})
	if err != nil {
		return err
	}

	fullname := getSchemaFileName(obj.appName, obj.dir, obj.name)
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

// Object List
type {{.PluralName}} []*{{.Name}}

`
