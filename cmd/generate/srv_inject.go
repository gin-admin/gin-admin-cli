package generate

import (
	"context"
	"fmt"
	"strings"
)

func getBllInjectFileName(appName, dir string) string {
	fullname := fmt.Sprintf("%s/internal/%s/service/service.go", dir, appName)
	return fullname
}

func insertServiceInject(ctx context.Context, obj *genObject) error {
	fullname := getBllInjectFileName(obj.appName, obj.dir)

	injectContent := fmt.Sprintf("%sSet,", obj.name)
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "var ServiceSet = wire.NewSet(") {
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

	err := insertContent(fullname, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}
