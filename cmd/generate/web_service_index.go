package generate

import (
	"context"
	"fmt"
	"strings"
)

func getWebServiceIndexFileName(dir string) string {
	fullname := fmt.Sprintf("%s/src/services/index.js", dir)
	return fullname
}

func insertWebServiceIndexImport(ctx context.Context, cmd *Command, item TplItem) error {
	name := strings.ToLower(item.StructName)
	injectContent := fmt.Sprintf("import %s from './%s';",
		item.StructName, name)
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if strings.Index(line, item.StructName) > -1 {
			injectStart = -1
		}
		if injectStart == 0 && strings.Contains(line, "export {") {
			injectStart = -1
			data = injectContent
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	filename := getWebServiceIndexFileName(cmd.cfg.Web)
	err := insertContent(filename, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", filename)

	return execGoFmt(filename)
}

func insertWebServiceIndexExport(ctx context.Context, cmd *Command, item TplItem) error {
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "export {") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Index(line, item.StructName) > -1 {
			injectStart = -1
		}
		if injectStart == 1 && strings.Contains(line, "};") {
			injectStart = -1
			data = "  " + item.StructName + ","
			flag = -1
			ok = true
			return
		}

		return "", 0, false
	}

	filename := getWebServiceIndexFileName(cmd.cfg.Web)
	err := insertContent(filename, insertFn)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", filename)

	return execGoFmt(filename)
}
