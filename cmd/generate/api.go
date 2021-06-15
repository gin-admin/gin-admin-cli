package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v4/util"
)

func getAPIFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/api/%s.api.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genAPI(ctx context.Context, pkgName, dir, name, comment string) error {
	pname := util.ToPlural(util.ToLowerUnderlinedNamer(name))
	pname = strings.Replace(pname, "_", "-", -1)

	data := map[string]interface{}{
		"PkgName": pkgName,
		"Name":    name,
		"Comment": comment,
	}

	buf, err := execParseTpl(apiTpl, data)
	if err != nil {
		return err
	}

	fullname := getAPIFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const apiTpl = `
package api

import (
	"{{.PkgName}}/internal/app/service"
	"{{.PkgName}}/internal/app/ginx"
	"{{.PkgName}}/internal/app/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// {{.Name}}Set 注入{{.Name}}
var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}), "*"))

// {{.Name}} {{.Comment}}
type {{.Name}} struct {
	{{.Name}}Srv *service.{{.Name}}
}

// Query 查询数据
func (a *{{.Name}}) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.{{.Name}}QueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.{{.Name}}Srv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}

	ginx.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
func (a *{{.Name}}) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.{{.Name}}Srv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// Create 创建数据
func (a *{{.Name}}) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = ginx.GetUserID(c)
	result, err := a.{{.Name}}Srv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// Update 更新数据
func (a *{{.Name}}) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.{{.Name}}Srv.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// Delete 删除数据
func (a *{{.Name}}) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

`
