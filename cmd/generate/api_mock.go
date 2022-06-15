package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getAPIMockFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/api/mock/%s.mock.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genAPIMock(ctx context.Context, pkgName, dir, name, comment string, excludeStatus, excludeCreate bool) error {
	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	pname = strings.Replace(pname, "_", "-", -1)

	data := map[string]interface{}{
		"PkgName":       pkgName,
		"Name":          name,
		"Comment":       comment,
		"PluralName":    util.ToPlural(pname),
		"IncludeStatus": !excludeStatus,
		"IncludeCreate": !excludeCreate,
	}

	buf, err := execParseTpl(apiMockTpl, data)
	if err != nil {
		return err
	}

	fullname := getAPIMockFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const apiMockTpl = `
package mock

import (
	_ "github.com/LyricTian/gin-admin/v8/internal/app/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}Mock), "*"))

type {{.Name}}Mock struct{}

// @Tags {{.Comment}}
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Success 200 {object} schema.ListResult{list=[]schema.{{.Name}}} "Response Data"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}} [get]
func (a *{{.Name}}Mock) Query(c *gin.Context) {
}

// @Tags {{.Comment}}
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.{{.Name}}
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:bad request}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}}/{id} [get]
func (a *{{.Name}}Mock) Get(c *gin.Context) {
}

// @Tags {{.Comment}}
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.{{.Name}} true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:bad request}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}} [post]
func (a *{{.Name}}Mock) Create(c *gin.Context) {
}

// @Tags {{.Comment}}
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Param body body schema.{{.Name}} true "更新数据"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:bad request}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}}/{id} [put]
func (a *{{.Name}}Mock) Update(c *gin.Context) {
}

// @Tags {{.Comment}}
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}}/{id} [delete]
func (a *{{.Name}}Mock) Delete(c *gin.Context) {
}

{{if .IncludeStatus}}
// @Tags {{.Comment}}
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}}/{id}/enable [patch]
func (a *{{.Name}}Mock) Enable(c *gin.Context) {
}

// @Tags {{.Comment}}
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:9999,message:invalid signature}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:internal server error}}"
// @Router /api/v1/{{.PluralName}}/{id}/disable [patch]
func (a *{{.Name}}Mock) Disable(c *gin.Context) {
}
{{end}}

`
