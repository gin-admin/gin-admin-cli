package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/LyricTian/gin-admin-cli/util"
)

func getAPIFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/api/a_%s.go", dir, util.ToLowerUnderlinedNamer(name))
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
	"{{.PkgName}}/internal/app/bll"
	"{{.PkgName}}/internal/app/ginplus"
	"{{.PkgName}}/internal/app/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// {{.Name}}Set 注入{{.Name}}
var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}), "*"))

// {{.Name}} {{.Comment}}
type {{.Name}} struct {
	{{.Name}}Bll bll.I{{.Name}}
}

// Query 查询数据
func (a *{{.Name}}) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.{{.Name}}QueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.{{.Name}}Bll.Query(ctx, params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get 查询指定数据
func (a *{{.Name}}) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.{{.Name}}Bll.Get(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create 创建数据
func (a *{{.Name}}) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	result, err := a.{{.Name}}Bll.Create(ctx, item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// Update 更新数据
func (a *{{.Name}}) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.{{.Name}}Bll.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Delete 删除数据
func (a *{{.Name}}) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Bll.Delete(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

`
