package generate

import (
	"context"
	"fmt"
	"strings"
)

func getAPIMockInjectFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/api/mock/mock.go", dir)
	return fullname
}

func insertAPIMockInject(ctx context.Context, dir, name string) error {
	injectContent := fmt.Sprintf("%sSet,", name)
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "var MockSet = wire.NewSet(") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, ")") {
			injectStart = -1
			data = injectContent
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	filename := getAPIMockInjectFileName(dir)
	err := insertContent(filename, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", filename)

	return execGoFmt(filename)
}
