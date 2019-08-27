package generate

import (
	"context"
	"fmt"
)

func getBllInjectFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/bll/impl/impl.go", dir)
	return fullname
}

// 插入bll注入文件
func insertBllInject(ctx context.Context, pkgName, dir, name, comment string) error {
	fullname := getBllInjectFileName(dir)

	err := insertFileContent(fullname, "func Inject", "container.Provide", fmt.Sprintf("container.Provide(internal.New%s, dig.As(new(bll.I%s)))\n", name, name))
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
