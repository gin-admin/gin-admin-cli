package generate

import (
	"context"
	"fmt"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getCtlFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/routers/api/ctl/c_%s.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成ctl文件
func genCtl(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName":    pkgName,
		"Name":       name,
		"PluralName": util.ToPlural(name),
		"Comment":    comment,
	}

	buf, err := execParseTpl(ctlTpl, data)
	if err != nil {
		return err
	}

	fullname := getBllFileName(dir, name)
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
	"{{.PkgName}}/internal/app/errors"
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

// {{.Name}} {{.Comment}}
// @Name {{.Name}}
// @Description {{.Comment}}控制器
type {{.Name}} struct {
	{{.Name}}Bll bll.I{{.Name}}
}

// Query 查询数据
func (a *{{.Name}}) Query(c *gin.Context) {
	switch c.Query("q") {
	case "page":
		a.QueryPage(c)
	default:
		ginplus.ResError(c, errors.ErrUnknownQuery)
	}
}

// QueryPage 查询分页数据
// @Summary 查询分页数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param current query int true "分页索引" 1
// @Param pageSize query int true "分页大小" 10
// @Success 200 []schema.{{.Name}} "查询结果：{list:列表数据,pagination:{current:页索引,pageSize:页大小,total:总数量}}"
// @Failure 400 schema.HTTPError "{error:{code:0,message:未知的查询类型}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/{{.PluralName}}?q=page
func (a *{{.Name}}) QueryPage(c *gin.Context) {
	var params schema.{{.Name}}QueryParam

	result, err := a.{{.Name}}Bll.Query(ginplus.NewContext(c), params, schema.{{.Name}}QueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
// @Summary 查询指定数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.{{.Name}}
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 404 schema.HTTPError "{error:{code:0,message:资源不存在}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router GET /api/v1/{{.PluralName}}/{id}
func (a *{{.Name}}) Get(c *gin.Context) {
	item, err := a.{{.Name}}Bll.Get(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create 创建数据
// @Summary 创建数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param body body schema.{{.Name}} true
// @Success 200 schema.{{.Name}}
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router POST /api/v1/{{.PluralName}}
func (a *{{.Name}}) Create(c *gin.Context) {
	var item schema.{{.Name}}
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	nitem, err := a.{{.Name}}Bll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update 更新数据
// @Summary 更新数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Param body body schema.{{.Name}} true
// @Success 200 schema.{{.Name}}
// @Failure 400 schema.HTTPError "{error:{code:0,message:无效的请求参数}}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router PUT /api/v1/{{.PluralName}}/{id}
func (a *{{.Name}}) Update(c *gin.Context) {
	var item schema.{{.Name}}
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.{{.Name}}Bll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete 删除数据
// @Summary 删除数据
// @Param Authorization header string false "Bearer 用户令牌"
// @Param id path string true "记录ID"
// @Success 200 schema.HTTPStatus "{status:OK}"
// @Failure 401 schema.HTTPError "{error:{code:0,message:未授权}}"
// @Failure 500 schema.HTTPError "{error:{code:0,message:服务器错误}}"
// @Router DELETE /api/v1/{{.PluralName}}/{id}
func (a *{{.Name}}) Delete(c *gin.Context) {
	err := a.{{.Name}}Bll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

`
