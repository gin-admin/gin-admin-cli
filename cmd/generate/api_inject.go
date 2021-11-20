package generate

import (
	"context"
	"fmt"
	"strings"
)

func getAPIInjectFileName(appName, dir string) string {
	fullname := fmt.Sprintf("%s/internal/%s/api/api.go", dir, appName)
	return fullname
}

func insertAPIInject(ctx context.Context, obj *genObject) error {
	injectContent := fmt.Sprintf("%sSet,", obj.name)
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "var APISet = wire.NewSet(") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, ") // end") {
			injectStart = -1
			data = injectContent
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	filename := getAPIInjectFileName(obj.appName, obj.dir)
	err := insertContent(filename, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", filename)

	return execGoFmt(filename)
}
