package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getAPIFileName(dir, routerName string) string {
	fullname := fmt.Sprintf("%s/internal/app/routers/%s/%s.go", dir, routerName, routerName)
	return fullname
}

// 插入api文件
func insertAPI(ctx context.Context, pkgName, dir, routerName, name, comment string) error {
	fullname := getAPIFileName(dir, routerName)

	injectContent := fmt.Sprintf("c%s *ctl.%s,", name, name)
	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	pname = strings.Replace(pname, "_", "-", -1)
	apiContent, err := execParseTpl(apiTpl, map[string]string{
		"Name":       name,
		"PluralName": pname,
	})
	if err != nil {
		return err
	}

	var injectStart, apiStart int
	apiStack := -1
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "container.Invoke") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, "error") {
			injectStart = -1
			data = injectContent
			flag = -1
			ok = true
			return
		}

		if apiStart == 0 && strings.Contains(line, "v1 := g.Group") {
			apiStart = 1
			return
		}

		if apiStart == 1 {
			if v := strings.TrimSpace(line); v == "{" {
				if apiStack == -1 {
					apiStack = 0
				}
				apiStack++
			} else if v == "}" {
				apiStack--
			}

			if apiStack == 0 {
				apiStack = -1
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

const apiTpl = `
// 注册/api/v1/{{.PluralName}}
g{{.Name}} := v1.Group("{{.PluralName}}")
{
	g{{.Name}}.GET("", c{{.Name}}.Query)
	g{{.Name}}.GET(":id", c{{.Name}}.Get)
	g{{.Name}}.POST("", c{{.Name}}.Create)
	g{{.Name}}.PUT(":id", c{{.Name}}.Update)
	g{{.Name}}.DELETE(":id", c{{.Name}}.Delete)
}

`
