package generate

import (
	"context"
	"fmt"
	"strings"
)

func getWebRouterFileName(dir string) string {
	fullname := fmt.Sprintf("%s/config/routes.js", dir)
	return fullname
}

func insertWebRouter(ctx context.Context, cmd *Command, item TplItem) error {
	name := strings.ToLower(item.StructName)
	injectContent := fmt.Sprintf("  { name: '%s', icon: 'icon-%s', path: '/%s', component: './%s', access: 'can%s' },",
		item.Comment, name, name, name, item.StructName)
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if strings.Index(line, item.StructName) > -1 {
			injectStart = -1
		}
		if injectStart == 0 && strings.Contains(line, "{ path: '/'") {
			injectStart = -1
			data = injectContent
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	filename := getWebRouterFileName(cmd.cfg.Web)
	err := insertContent(filename, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", filename)

	return execGoFmt(filename)
}
