package generate

import (
	"context"
	"fmt"
)

func getModelInjectFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/gorm/gorm.go", dir)
	return fullname
}

// 插入model注入文件
func insertModelInject(ctx context.Context, pkgName, dir, name, comment string) error {
	fullname := getModelInjectFileName(dir)

	err := insertFileContent(fullname, "func AutoMigrate", "entity.", fmt.Sprintf("new(entity.%s),\n", name))
	if err != nil {
		return err
	}

	err = insertFileContent(fullname, "func Inject", "container.Provide", fmt.Sprintf("container.Provide(imodel.New%s, dig.As(new(model.I%s)))\n", name, name))
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
