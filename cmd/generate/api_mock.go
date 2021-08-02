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

func genAPIMock(ctx context.Context, pkgName, dir, name, comment string) error {
	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	pname = strings.Replace(pname, "_", "-", -1)

	data := map[string]interface{}{
		"PkgName":    pkgName,
		"Name":       name,
		"Comment":    comment,
		"PluralName": util.ToPlural(pname),
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
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}Mock), "*"))

// {{.Name}}Mock {{.Comment}}
type {{.Name}}Mock struct{}

// Query 查询数据
// @Tags {{.Comment}}
// @Summary 查询数据
// @Security ApiKeyAuth
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Param queryValue query string false "查询值"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Success 200 {object} schema.ListResult{list=[]schema.{{.Name}}} "查询结果"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}} [get]
func (a *{{.Name}}Mock) Query(c *gin.Context) {
}

// Get 查询指定数据
// @Tags {{.Comment}}
// @Summary 查询指定数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.{{.Name}}
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.ErrorResult "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id} [get]
func (a *{{.Name}}Mock) Get(c *gin.Context) {
}

// Create 创建数据
// @Tags {{.Comment}}
// @Summary 创建数据
// @Security ApiKeyAuth
// @Param body body schema.{{.Name}} true "创建数据"
// @Success 200 {object} schema.IDResult
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}} [post]
func (a *{{.Name}}Mock) Create(c *gin.Context) {
}

// Update 更新数据
// @Tags {{.Comment}}
// @Summary 更新数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Param body body schema.{{.Name}} true "更新数据"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 400 {object} schema.ErrorResult "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id} [put]
func (a *{{.Name}}Mock) Update(c *gin.Context) {
}

// Delete 删除数据
// @Tags {{.Comment}}
// @Summary 删除数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id} [delete]
func (a *{{.Name}}Mock) Delete(c *gin.Context) {
}

// Enable 启用数据
// @Tags {{.Comment}}
// @Summary 启用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id}/enable [patch]
func (a *{{.Name}}Mock) Enable(c *gin.Context) {
}

// Disable 禁用数据
// @Tags {{.Comment}}
// @Summary 禁用数据
// @Security ApiKeyAuth
// @Param id path int true "唯一标识"
// @Success 200 {object} schema.StatusResult "{status:OK}"
// @Failure 401 {object} schema.ErrorResult "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.ErrorResult "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id}/disable [patch]
func (a *{{.Name}}Mock) Disable(c *gin.Context) {
}

`
