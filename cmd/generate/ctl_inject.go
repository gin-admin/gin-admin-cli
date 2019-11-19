package generate

import (
	"context"
	"fmt"
	"strings"
)

func getCtlInjectFileName(dir, routerName string) string {
	fullname := fmt.Sprintf("%s/internal/app/routers/%s/ctl/ctl.go", dir, routerName)
	return fullname
}

// 插入ctl注入文件
func insertCtlInject(ctx context.Context, pkgName, dir, routerName, name, comment string) error {
	fullname := getCtlInjectFileName(dir, routerName)

	injectContent := fmt.Sprintf("_ = container.Provide(New%s)", name)
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "container *dig.Container") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, "return") {
			injectStart = -1
			data = injectContent
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	err := insertContent(fullname, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
