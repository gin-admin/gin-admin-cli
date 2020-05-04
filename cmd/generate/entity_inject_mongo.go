package generate

import (
	"context"
	"fmt"
	"strings"
)

func getEntityInjectMongoFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/mongo/mongo.go", dir)
	return fullname
}

func insertEntityInjectMongo(ctx context.Context, dir, name string) error {
	fullname := getEntityInjectMongoFileName(dir)

	injectContent := fmt.Sprintf("new(entity.%s),", name)

	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "return createIndexes(") {
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

	err := insertContent(fullname, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
