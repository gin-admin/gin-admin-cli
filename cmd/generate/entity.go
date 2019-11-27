package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin-cli/util"
)

type entityField struct {
	Name        string // 字段名
	Comment     string // 字段注释
	Type        string // 字段类型
	GormOptions string // gorm配置项(不包含column)
}

func getEntityFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/gorm/internal/entity/e_%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成entity文件
func genEntity(ctx context.Context, pkgName, dir, name, comment string, fields ...entityField) error {
	if len(fields) == 0 {
		fields = []entityField{
			{Name: "RecordID", Comment: "记录ID", Type: "string", GormOptions: "size:36;index;"},
			{Name: "Creator", Comment: "创建者", Type: "string", GormOptions: "size:36;index;"},
		}
	}

	buf := new(bytes.Buffer)

	var imports []string
	imports = append(imports, `"context"`)
	imports = append(imports, fmt.Sprintf(`"%s/internal/app/schema"`, pkgName))
	imports = append(imports, `"github.com/jinzhu/gorm"`)

	buf.Write(getModuleHeader("entity", imports...).Bytes())

	buf.WriteString(fmt.Sprintf("// Get%sDB %s", name, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("func Get%sDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {", name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("return getDBWithModel(ctx, defDB, %s{})", name))
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("// Schema%s %s", name, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("type Schema%s schema.%s", name, name))
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("// To%s 转换为%s实体", name, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("func (a Schema%s) To%s() *%s {", name, name, name))
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("item := &%s{", name))
	buf.WriteString(delimiter)

	for _, field := range fields {
		buf.WriteString(fmt.Sprintf("%s: &a.%s,", field.Name, field.Name))
		buf.WriteString(delimiter)
	}

	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("return item"))
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("// %s %s实体", name, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("type %s struct {", name))
	buf.WriteString(delimiter)
	buf.WriteString("Model")
	buf.WriteString(delimiter)
	for _, field := range fields {
		buf.WriteString(fmt.Sprintf("%s *%s", field.Name, field.Type))
		buf.WriteByte('`')
		buf.WriteString(fmt.Sprintf(`gorm:"column:%s;%s"`, util.ToLowerUnderlinedNamer(field.Name), field.GormOptions))
		buf.WriteByte('`')
		buf.WriteString(fmt.Sprintf(" // %s", field.Comment))
		buf.WriteString(delimiter)
	}
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("func (a %s) String() string {", name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("return toString(a)"))
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString("// TableName 表名")
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("func (a %s) TableName() string {", name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`return a.Model.TableName("%s")`, util.ToLowerUnderlinedNamer(name)))
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("// ToSchema%s 转换为%s对象", name, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("func (a %s) ToSchema%s() *schema.%s {", name, name, name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("item := &schema.%s{", name))
	buf.WriteString(delimiter)

	for _, field := range fields {
		buf.WriteString(fmt.Sprintf("%s:  *a.%s,", field.Name, field.Name))
		buf.WriteString(delimiter)
	}
	buf.WriteString(fmt.Sprintf("CreatedAt:  a.CreatedAt,"))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("UpdatedAt:  a.UpdatedAt,"))
	buf.WriteString(delimiter)

	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString("return item")
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	pluralName := util.ToPlural(name)
	buf.WriteString(fmt.Sprintf("// %s %s列表", pluralName, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("type %s []*%s", pluralName, name))
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("// ToSchema%s 转换为%s对象列表", pluralName, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("func (a %s) ToSchema%s() []*schema.%s {", pluralName, pluralName, name))
	buf.WriteString(delimiter)

	buf.WriteString(fmt.Sprintf("list := make([]*schema.%s, len(a))", name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("for i, item := range a {"))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("list[i] = item.ToSchema%s()", name))
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString("return list")
	buf.WriteString(delimiter)
	buf.WriteByte('}')
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	fullname := getEntityFileName(dir, name)
	err := createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
