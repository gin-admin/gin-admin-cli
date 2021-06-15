package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v4/util"
)

type entityMongoField struct {
	Name    string // 字段名
	Comment string // 字段注释
	Type    string // 字段类型
}

func getEntityMongoFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/mongo/entity/%s.entity.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成entity文件
func genMongoEntity(ctx context.Context, pkgName, dir, name, comment string, fields ...entityMongoField) error {
	var tfields []entityMongoField

	tfields = append(tfields, entityMongoField{Name: "Model", Comment: "", Type: ""})
	tfields = append(tfields, fields...)
	tfields = append(tfields, entityMongoField{Name: "Creator", Comment: "创建者", Type: "string"})

	buf := new(bytes.Buffer)
	for _, field := range tfields {
		buf.WriteString(fmt.Sprintf("%s \t %s \t", field.Name, field.Type))
		buf.WriteByte('`')
		if field.Type == "" {
			buf.WriteString(`bson:",inline"`)
		} else {
			buf.WriteString(fmt.Sprintf(`bson:"%s"`, util.ToLowerUnderlinedNamer(field.Name)))
		}
		buf.WriteByte('`')

		if field.Comment != "" {
			buf.WriteString(fmt.Sprintf("// %s", field.Comment))
		}
		buf.WriteString(delimiter)
	}

	tbuf, err := execParseTpl(entityMongoTpl, map[string]interface{}{
		"PkgName":       pkgName,
		"Name":          name,
		"PluralName":    util.ToPlural(name),
		"Fields":        buf.String(),
		"Comment":       comment,
		"UnderLineName": util.ToLowerUnderlinedNamer(name),
	})
	if err != nil {
		return err
	}

	fullname := getEntityMongoFileName(dir, name)
	err = createFile(ctx, fullname, tbuf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const entityMongoTpl = `
package entity

import (
	"context"

	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Get{{.Name}}Collection 获取{{.Name}}集合
func Get{{.Name}}Collection(ctx context.Context, cli *mongo.Client) *mongo.Collection {
	return GetCollection(ctx, cli, {{.Name}}{})
}

// Schema{{.Name}} {{.Comment}}
type Schema{{.Name}} schema.{{.Name}}

// To{{.Name}} 转换为实体
func (a Schema{{.Name}}) To{{.Name}}() *{{.Name}} {
	item := new({{.Name}})
	util.StructMapToStruct(a, item)
	return item
}

// {{.Name}} {{.Comment}}实体
type {{.Name}} struct {
	{{.Fields}}
}

// CollectionName 集合名
func (a {{.Name}}) CollectionName() string {
	return a.Model.CollectionName("{{.UnderLineName}}")
}

// CreateIndexes 创建索引
func (a {{.Name}}) CreateIndexes(ctx context.Context, cli *mongo.Client) error {
	return a.Model.CreateIndexes(ctx, cli, a, []mongo.IndexModel{
		{Keys: bson.M{"creator": 1}},
	})
}

// ToSchema{{.Name}} 转换为对象
func (a {{.Name}}) ToSchema{{.Name}}() *schema.{{.Name}} {
	item := new(schema.{{.Name}})
	util.StructMapToStruct(a, item)
	return item
}

// {{.PluralName}} {{.Comment}}实体列表
type {{.PluralName}} []*{{.Name}}

// ToSchema{{.Name}}s 转换为{{.Comment}}对象列表
func (a {{.PluralName}}) ToSchema{{.PluralName}}() schema.{{.PluralName}} {
	list := make(schema.{{.PluralName}}, len(a))
	for i, item := range a {
		list[i] = item.ToSchema{{.Name}}()
	}
	return list
}

`
