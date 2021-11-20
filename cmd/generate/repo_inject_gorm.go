package generate

import (
	"context"
	"fmt"
	"strings"
)

func getModelInjectGormFileName(appName, dir string) string {
	fullname := fmt.Sprintf("%s/internal/%s/dao/dao.go", dir, appName)
	return fullname
}

func insertRepoInjectGorm(ctx context.Context, obj *genObject) error {
	fullname := getModelInjectGormFileName(obj.appName, obj.dir)

	injectStart2 := 0
	injectStart3 := 0
	injectStart4 := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart2 == 0 && strings.Contains(line, "var RepoSet = wire.NewSet(") {
			injectStart2 = 1
			return
		}

		if injectStart2 == 1 && strings.Contains(line, ") // end") {
			injectStart2 = -1
			data = "\t" + fmt.Sprintf(`repo.%sSet,`, obj.name)
			flag = -1
			ok = true
			return
		}

		if injectStart3 == 0 && strings.Contains(line, "type (") {
			injectStart3 = 1
			return
		}

		if injectStart3 == 1 && strings.Contains(line, ") // end") {
			injectStart3 = -1
			data = "\t" + fmt.Sprintf(`%sRepo               = repo.%sRepo`, obj.name, obj.name)
			flag = -1
			ok = true
			return
		}

		if injectStart4 == 0 && strings.Contains(line, "func AutoMigrate(db *gorm.DB) error {") {
			injectStart4 = 1
			return
		}

		if injectStart4 == 1 && strings.Contains(line, ") // end") {
			injectStart4 = -1
			data = "\t" + fmt.Sprintf(`new(schema.%s),`, obj.name)
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
