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
func insertBllInject(ctx context.Context, dir, name string) error {
	fullname := getBllInjectFileName(dir)

	err := insertFileContent(fullname, "func Inject", "container.Provide", fmt.Sprintf("container.Provide(internal.New%s, dig.As(new(bll.I%s)))\n", name, name))
	if err != nil {
		return err
	}

	return execGoFmt(fullname)
}
