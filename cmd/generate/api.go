package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getAPIFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/api/%s.api.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

func genAPI(ctx context.Context, pkgName, dir, name, comment string) error {
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

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const apiTpl = `
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"{{.PkgName}}/internal/app/contextx"
	"{{.PkgName}}/internal/app/ginx"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/internal/app/service"
)

var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}API), "*"))

type {{.Name}}API struct {
	{{.Name}}Srv *service.{{.Name}}Srv
}

func (a *{{.Name}}API) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.{{.Name}}QueryParam
	if err := ginx.ParseQuery(c, &params); err != nil {
		ginx.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.{{.Name}}Srv.Query(ctx, params, schema.{{.Name}}QueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResPage(c, result.Data, result.PageResult)
}

func (a *{{.Name}}API) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.{{.Name}}Srv.Get(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, item)
}

func (a *{{.Name}}API) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	item.Creator = contextx.FromUserID(ctx)
	result, err := a.{{.Name}}Srv.Create(ctx, item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResSuccess(c, result)
}

func (a *{{.Name}}API) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.{{.Name}}
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, err)
		return
	}

	err := a.{{.Name}}Srv.Update(ctx, ginx.ParseParamID(c, "id"), item)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *{{.Name}}API) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.Delete(ctx, ginx.ParseParamID(c, "id"))
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *{{.Name}}API) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 1)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

func (a *{{.Name}}API) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.{{.Name}}Srv.UpdateStatus(ctx, ginx.ParseParamID(c, "id"), 2)
	if err != nil {
		ginx.ResError(c, err)
		return
	}
	ginx.ResOK(c)
}

`
