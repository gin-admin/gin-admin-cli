package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

type entityGormField struct {
	Name        string // 字段名
	Comment     string // 字段注释
	Type        string // 字段类型
	GormOptions string // gorm配置项(不包含column)
}

func getEntityGormFileName(dir, name string) string {
	name = util.ToLowerUnderlinedNamer(name)
	fullname := fmt.Sprintf("%s/internal/app/dao/%s/%s.entity.go", dir, name, name)
	return fullname
}

func genGormEntity(ctx context.Context, pkgName, dir, name, comment string, excludeStatus, excludeCreate bool, fields ...entityGormField) error {
	var tfields []entityGormField

	tfields = append(tfields, fields...)
	if !excludeStatus {
		tfields = append(tfields, entityGormField{Name: "Status", Comment: "状态(1:启用 2:停用)", Type: "int", GormOptions: "type:tinyint;index;default:0;"})
	}

	if !excludeCreate {
		tfields = append(tfields, entityGormField{Name: "Creator", Comment: "创建者", Type: "uint64"})
	}

	buf := new(bytes.Buffer)
	for _, field := range tfields {
		buf.WriteString(fmt.Sprintf("%s \t %s \t", field.Name, field.Type))
		buf.WriteByte('`')

		gormTag := ""
		if field.GormOptions != "" {
			gormTag += field.GormOptions
		}
		buf.WriteString(fmt.Sprintf(`gorm:"%s"`, gormTag))

		buf.WriteByte('`')

		if field.Comment != "" {
			buf.WriteString(fmt.Sprintf("// %s", field.Comment))
		}
		buf.WriteString(delimiter)
	}

	tbuf, err := execParseTpl(entityGormTpl, map[string]interface{}{
		"PkgName":       pkgName,
		"Name":          name,
		"PluralName":    util.ToPlural(name),
		"Fields":        buf.String(),
		"Comment":       comment,
		"UnderLineName": util.ToLowerUnderlinedNamer(name),
		"BackQuote":     "`",
	})
	if err != nil {
		return err
	}

	filename := getEntityGormFileName(dir, name)
	err = createFile(ctx, filename, tbuf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", filename)

	return execGoFmt(filename)
}

const entityGormTpl = `
package {{.UnderLineName}}

import (
	"context"

	"gorm.io/gorm"

	"{{.PkgName}}/internal/app/dao/util"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/util/structure"
)

// Get {{.Name}} db model
func Get{{.Name}}DB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new({{.Name}}))
}

// {{.Name}}
type Schema{{.Name}} schema.{{.Name}}

// Convert to {{.Name}} entity
func (a Schema{{.Name}}) To{{.Name}}() *{{.Name}} {
	item := new({{.Name}})
	structure.Copy(a, item)
	return item
}

// {{.Name}} entity
type {{.Name}} struct {
	util.Model
	{{.Fields}}
}

// Convert to {{.Name}} schema
func (a {{.Name}}) ToSchema{{.Name}}() *schema.{{.Name}} {
	item := new(schema.{{.Name}})
	structure.Copy(a, item)
	return item
}

// {{.Name}} entity list
type {{.PluralName}} []*{{.Name}}

// Convert to {{.Name}} schema list
func (a {{.PluralName}}) ToSchema{{.PluralName}}() []*schema.{{.Name}} {
	list := make([]*schema.{{.Name}}, len(a))
	for i, item := range a {
		list[i] = item.ToSchema{{.Name}}()
	}
	return list
}
`
