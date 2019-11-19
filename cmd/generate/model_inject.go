package generate

import (
	"context"
	"fmt"
	"strings"
)

func getModelInjectFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/gorm/gorm.go", dir)
	return fullname
}

// 插入model注入文件
func insertModelInject(ctx context.Context, pkgName, dir, name, comment string) error {
	fullname := getModelInjectFileName(dir)

	migrateContent := fmt.Sprintf("new(entity.%s),", name)
	injectContent := fmt.Sprintf("_ = container.Provide(imodel.New%s, dig.As(new(model.I%s)))", name, name)

	migrateStart := 0
	injectStart := 0
	insertFn := func(line string) (data string, flag int, ok bool) {
		if migrateStart == 0 && strings.Contains(line, "db.AutoMigrate") {
			migrateStart = 1
			return
		}

		if migrateStart == 1 && strings.Contains(line, "Error") {
			migrateStart = -1
			data = migrateContent
			flag = -1
			ok = true
			return
		}

		if injectStart == 0 && strings.Contains(line, "container *dig.Container") {
			injectStart = 1
			return
		}

		if injectStart == 1 && strings.Contains(line, "return") {
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
