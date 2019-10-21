package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getCtlFileName(dir, routerName, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/routers/%s/ctl/c_%s.go", dir, routerName, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成ctl文件
func genCtl(ctx context.Context, pkgName, dir, routerName, name, comment string) error {
	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	pname = strings.Replace(pname, "_", "-", -1)

	data := map[string]interface{}{
		"PkgName":    pkgName,
		"Name":       name,
		"PluralName": pname,
		"Comment":    comment,
	}

	buf, err := execParseTpl(ctlTpl, data)
	if err != nil {
		return err
	}

	fullname := getCtlFileName(dir, routerName, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const ctlTpl = `
package ctl

import (
	"{{.PkgName}}/internal/app/bll"
	"{{.PkgName}}/internal/app/ginplus"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/util"
	"github.com/gin-gonic/gin"
)

// New{{.Name}} 创建{{.Comment}}控制器
func New{{.Name}}(b{{.Name}} bll.I{{.Name}}) *{{.Name}} {
	return &{{.Name}}{
		{{.Name}}Bll: b{{.Name}},
	}
}

// {{.Name}} {{.Comment}}控制器
type {{.Name}} struct {
	{{.Name}}Bll bll.I{{.Name}}
}

// Query 查询数据
// @Tags {{.Comment}}
// @Summary 查询数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" default(1)
// @Param pageSize query int true "分页大小" default(10)
// @Success 200 {array} schema.{{.Name}} "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}} [get]
func (a *Demo) Query(c *gin.Context) {
	var params schema.{{.Name}}QueryParam

	result, err := a.DemoBll.Query(ginplus.NewContext(c), params, schema.{{.Name}}QueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// @Tags {{.Comment}}
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.{{.Name}}
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id} [get]
func (a *Demo) Get(c *gin.Context) {
	item, err := a.DemoBll.Get(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create 创建数据
// @Tags {{.Comment}}
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.{{.Name}} true "创建数据"
// @Success 200 {object} schema.{{.Name}}
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}} [post]
func (a *Demo) Create(c *gin.Context) {
	var item schema.{{.Name}}
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	nitem, err := a.DemoBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update 更新数据
// @Tags {{.Comment}}
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.{{.Name}} true "更新数据"
// @Success 200 {object} schema.{{.Name}}
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id} [put]
func (a *Demo) Update(c *gin.Context) {
	var item schema.{{.Name}}
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.DemoBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete 删除数据
// @Tags {{.Comment}}
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router /api/v1/{{.PluralName}}/{id} [delete]
func (a *Demo) Delete(c *gin.Context) {
	err := a.DemoBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

`
