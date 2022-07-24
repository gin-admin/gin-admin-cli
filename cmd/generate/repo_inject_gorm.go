package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getModelInjectGormFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/dao/dao.go", dir)
	return fullname
}

func insertModelInjectGorm(ctx context.Context, pkgName, dir, name string) error {
	fullname := getModelInjectGormFileName(dir)
	ulname := util.ToLowerUnderlinedNamer(name)

	injectStart := 0
	injectStart2 := 0
	injectStart3 := 0
	injectStart4 := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if injectStart == 0 && strings.Contains(line, "import (") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, ") // end") {
			injectStart = -1
			data = "\t" + fmt.Sprintf(`"%s/internal/app/dao/%s"`, pkgName, ulname)
			flag = -1
			ok = true
			return
		}

		if injectStart2 == 0 && strings.Contains(line, "var RepoSet = wire.NewSet(") {
			injectStart2 = 1
			return
		}

		if injectStart2 == 1 && strings.Contains(line, ") // end") {
			injectStart2 = -1
			data = "\t" + fmt.Sprintf(`%s.%sSet,`, ulname, name)
			flag = -1
			ok = true
			return
		}

		if injectStart3 == 0 && strings.Contains(line, "type (") {
			injectStart3 = 1
			return
		}

		if injectStart3 == 1 && strings.Index(line, name) > -1 {
			injectStart3 = -1
		}
		if injectStart3 == 1 && strings.Contains(line, ") // end") {
			injectStart3 = -1
			data = "\t" + fmt.Sprintf(`%sRepo               = %s.%sRepo`, name, ulname, name)
			flag = -1
			ok = true
			return
		}

		if injectStart4 == 0 && strings.Contains(line, "func AutoMigrate(db *gorm.DB) error {") {
			injectStart4 = 1
			return
		}

		if injectStart4 == 1 && strings.Index(line, name) > -1 {
			injectStart4 = -1
		}
		if injectStart4 == 1 && strings.Contains(line, ") // end") {
			injectStart4 = -1
			data = "\t" + fmt.Sprintf(`new(%s.%s),`, ulname, name)
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
