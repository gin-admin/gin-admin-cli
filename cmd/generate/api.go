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
func insertAPI(ctx context.Context, dir, name string) error {
	fullname := getAPIFileName(dir)

	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))

	buf := new(bytes.Buffer)
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf("// 注册/api/v1/%s", pname))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.GET("/%s", user.Query)`, pname))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.GET("/%s/:id", user.Get)`, pname))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.POST("/%s", user.Create)`, pname))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.PUT("/%s/:id", user.Update)`, pname))
	buf.WriteString(delimiter)
	buf.WriteString(fmt.Sprintf(`v1.DELETE("/%s/:id", user.Delete)`, pname))

	err := insertFileContent(fullname, "v1 := g.Group", "v1.", buf.String(), "//", "pub")
	if err != nil {
		return err
	}

	return execGoFmt(fullname)
}
