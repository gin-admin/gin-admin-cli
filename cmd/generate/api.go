package generate

import (
	"bytes"
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getAPIFileName(dir string) string {
	fullname := fmt.Sprintf("%s/internal/app/routers/api/api.go", dir)
	return fullname
}

// 插入api文件
func insertAPI(ctx context.Context, pkgName, dir, name, comment string) error {
	fullname := getAPIFileName(dir)

	err := insertFileContent(fullname, "return container.Invoke", "*ctl.", fmt.Sprintf("c%s *ctl.%s,\n", name, name))
	if err != nil {
		return err
	}

	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	buf := new(bytes.Buffer)
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("// 注册/api/v1/%s", pname))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.GET("/%s", c%s.Query)`, pname, name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.GET("/%s/:id", c%s.Get)`, pname, name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.POST("/%s", c%s.Create)`, pname, name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.PUT("/%s/:id", c%s.Update)`, pname, name))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.DELETE("/%s/:id", c%s.Delete)`, pname, name))

	err = insertFileContent(fullname, "v1 := g.Group", "v1.", buf.String(), "//", "pub")
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}
