package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getRouterAPIFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/router/r_api.go", dir)
	return fullname
}

func insertRouterAPI(ctx context.Context, dir, name string) error {
	fullname := getRouterAPIFileName(dir)

	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	pname = strings.Replace(pname, "_", "-", -1)
	apiContent, err := execParseTpl(routerAPITpl, map[string]string{
		"Name":       name,
		"PluralName": pname,
	})
	if err != nil {
		return err
	}

	var apiStart int
	apiStack := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if apiStart == 0 && strings.Contains(line, "v1 := g.Group") {
			apiStart = 1
			return
		}

		if apiStart == 1 {
			if v := strings.TrimSpace(line); v == "{" {
				apiStack++
			} else if v == "}" {
				apiStack--
			}

			if apiStack == 0 {
				data = apiContent.String()
				flag = -1
				ok = true
				return
			}
		}

		return "", 0, false
	}

	err = insertContent(fullname, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const routerAPITpl = `

g{{.Name}} := v1.Group("{{.PluralName}}")
{
	g{{.Name}}.GET("", a.{{.Name}}API.Query)
	g{{.Name}}.GET(":id", a.{{.Name}}API.Get)
	g{{.Name}}.POST("", a.{{.Name}}API.Create)
	g{{.Name}}.PUT(":id", a.{{.Name}}API.Update)
	g{{.Name}}.DELETE(":id", a.{{.Name}}API.Delete)
}
`
