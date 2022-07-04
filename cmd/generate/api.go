package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v6/util"
)

func getAPIFileName(appName, dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/%s/api/%s.api.go", dir, appName, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genAPI(ctx context.Context, obj *genObject) error {
	pname := util.ToPlural(util.ToLowerUnderlinedNamer(obj.name))
	data := map[string]interface{}{
		"PkgName":       obj.pkgName,
		"AppName":       obj.appName,
		"Name":          obj.name,
		"LowerName":     strings.Replace(pname, "_", " ", -1),
		"PluralName":    strings.Replace(pname, "_", "", -1),
		"Comment":       obj.comment,
		"IncludeStatus": !obj.excludeStatus,
		"IncludeCreate": !obj.excludeCreate,
	}

	buf, err := execParseTpl(apiTpl, data)
	if err != nil {
		return err
	}

	fullname := getAPIFileName(obj.appName, obj.dir, obj.name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const apiTpl = `
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	{{if .IncludeCreate}}
	"{{.PkgName}}/internal/{{.AppName}}/contextx"
	{{end}}
	"{{.PkgName}}/internal/{{.AppName}}/module/ginx"
	"{{.PkgName}}/internal/{{.AppName}}/schema"
	"{{.PkgName}}/internal/{{.AppName}}/service"
)

var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}API), "*"))

type {{.Name}}API struct {
	{{.Name}}Srv *service.{{.Name}}Srv
}

// @Tags {{.Name}}API
// @Summary Query {{.LowerName}} list
// @Security ApiKeyAuth
// @Param current query int true "pagination index" default(1)
// @Param pageSize query int true "pagination size" default(10)
// @Success 200 {object} schema.ListResult{list=[]schema.{{.Name}}} "Query result (schema.{{.Name}} object)"
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}} [get]
func (a *{{.Name}}API) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.{{.Name}}QueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.{{.Name}}Srv.Query(ctx, params)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

// @Tags {{.Name}}API
// @Summary Get single {{.LowerName}} by id
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Success 200 {object} schema.{{.Name}}
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}}/{id} [get]
func (a *{{.Name}}API) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.{{.Name}}Srv.Get(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

// @Tags {{.Name}}API
// @Summary Create {{.LowerName}}
// @Security ApiKeyAuth
// @Param body body schema.{{.Name}} true "Request body"
// @Success 200 {object} schema.{{.Name}}
// @Failure 400 {object} schema.ErrorResult
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}} [post]
func (a *{{.Name}}API) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	result, err := a.{{.Name}}Srv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

// @Tags {{.Name}}API
// @Summary Update {{.LowerName}} by id
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Param body body schema.{{.Name}} true "Request body"
// @Success 200 {object} schema.OkResult "ok=true"
// @Failure 400 {object} schema.ErrorResult
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}}/{id} [put]
func (a *{{.Name}}API) Update(c *gin.Context) {
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

// @Tags {{.Name}}API
// @Summary Delete single {{.LowerName}} by id
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Success 200 {object} schema.OkResult "ok=true"
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}}/{id} [delete]
func (a *{{.Name}}API) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.Delete(ctx, c.Param("id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

{{if .IncludeStatus}}
// @Tags {{.Name}}API
// @Summary Set {{.LowerName}} to enable
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Success 200 {object} schema.OkResult "ok=true"
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}}/{id}/enable [patch]
func (a *{{.Name}}API) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

// @Tags {{.Name}}API
// @Summary Set {{.LowerName}} to disable
// @Security ApiKeyAuth
// @Param id path string true "unique id"
// @Success 200 {object} schema.OkResult "ok=true"
// @Failure 401 {object} schema.ErrorResult
// @Failure 500 {object} schema.ErrorResult
// @Router /api/v1/{{.PluralName}}/{id}/disable [patch]
func (a *{{.Name}}API) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}
{{end}}

`
