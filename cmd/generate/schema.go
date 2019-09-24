package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin-cli/util"
)

type schemaField struct {
	Name       string // 字段名
	Comment    string // 字段注释
	Type       string // 字段类型
	IsRequired bool   // 是否必选项
}

func getSchemaFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/schema/s_%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成schema文件
func genSchema(ctx context.Context, dir, name, comment string, fields ...schemaField) error {
	if len(fields) == 0 {
		fields = []schemaField{
			{Name: "RecordID", Comment: "记录ID", Type: "string"},
			{Name: "Creator", Comment: "创建者", Type: "string"},
		}
	}

	buf := new(bytes.Buffer)

	buf.Write(getModuleHeader("schema").Bytes())

	buf.WriteString(fmt.Sprintf("// %s %s", name, comment))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("type %s struct {", name))
	buf.WriteString(delimiter)

	for _, field := range fields {
		buf.WriteString(fmt.Sprintf("%s \t %s \t", field.Name, field.Type))
		buf.WriteByte('`')
		buf.WriteString(fmt.Sprintf(`json:"%s"`, util.ToLowerUnderlinedNamer(field.Name)))

		if field.IsRequired {
			buf.WriteByte(' ')
			buf.WriteString(fmt.Sprintf(`binding:"required"`))
		}

		swagRequired := "false"
		if field.IsRequired {
			swagRequired = "true"
		}
		buf.WriteByte(' ')
		buf.WriteString(fmt.Sprintf(`swaggo:"%s,%s"`, swagRequired, field.Comment))

		buf.WriteByte('`')
		buf.WriteString(delimiter)
	}

	buf.WriteString("}")
	buf.WriteString(delimiter)

	tbuf, err := execParseTpl(schemaTpl, map[string]interface{}{
		"Name":       name,
		"PluralName": util.ToPlural(name),
		"Comment":    comment,
	})
	if err != nil {
		return err
	}

	buf.Write(tbuf.Bytes())

	fullname := getSchemaFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const schemaTpl = `
// {{.Name}}QueryParam 查询条件
type {{.Name}}QueryParam struct {
}

// {{.Name}}QueryOptions 查询可选参数项
type {{.Name}}QueryOptions struct {
	PageParam *PaginationParam // 分页参数
}

// {{.Name}}QueryResult 查询结果
type {{.Name}}QueryResult struct {
	Data       {{.PluralName}}
	PageResult *PaginationResult
}

// {{.PluralName}} {{.Comment}}列表
type {{.PluralName}} []*{{.Name}}

`
