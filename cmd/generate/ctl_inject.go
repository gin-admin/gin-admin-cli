package generate

import (
	"context"
	"fmt"
)

func getCtlInjectFileName(dir, routerName string) string {
	fullname := fmt.Sprintf("%s/internal/app/routers/%s/ctl/ctl.go", dir, routerName)
	return fullname
}

// 插入ctl注入文件
func insertCtlInject(ctx context.Context, pkgName, dir, routerName, name, comment string) error {
	fullname := getCtlInjectFileName(dir, routerName)

	err := insertFileContent(fullname, "func Inject", "container.Provide", fmt.Sprintf("container.Provide(New%s)\n", name))
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
