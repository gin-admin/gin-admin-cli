package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v6/util"
)

func getRouterAPIFileName(appName, dir string) string {
	fullname := fmt.Sprintf("%s/internal/%s/router/router.go", dir, appName)
	return fullname
}

func insertRouterAPI(ctx context.Context, obj *genObject) error {
	fullname := getRouterAPIFileName(obj.appName, obj.dir)

	pname := util.ToPlural(util.ToLowerUnderlinedNamer(obj.name))
	pname = strings.Replace(pname, "_", "-", -1)
	injectContent, err := execParseTpl(routerAPITpl, map[string]interface{}{
		"Name":          obj.name,
		"PluralName":    pname,
		"LowerName":     strings.ToLower(obj.name),
		"IncludeStatus": !obj.excludeStatus,
		"IncludeCreate": !obj.excludeCreate,
	})
	if err != nil {
		return err
	}

	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "v1 := g.Group") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, "} // v1 end") {
			injectStart = -1
			data = injectContent.String()
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	err = insertContent(fullname, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

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
	{{if .IncludeStatus}}
	g{{.Name}}.PATCH(":id/enable", a.{{.Name}}API.Enable)
	g{{.Name}}.PATCH(":id/disable", a.{{.Name}}API.Disable)
	{{end}}
}
`
