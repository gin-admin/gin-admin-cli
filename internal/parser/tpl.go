package parser

var tplModuleMain = `
package $$LowerModuleName$$

import (
	"context"
	"$$RootImportPath$$/internal/config"

	"$$ModuleImportPath$$/api"
	"$$ModuleImportPath$$/schema"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type $$ModuleName$$ struct {
	DB          *gorm.DB
}

func (a *$$ModuleName$$) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate()
}

func (a *$$ModuleName$$) Init(ctx context.Context) error {
	if config.C.Storage.DB.AutoMigrate {
		if err := a.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *$$ModuleName$$) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	{{- if .FillRouterPrefix}}
	v1 = v1.Group("$$LowerModuleName$$")
	{{- end}}
	return nil
}

func (a *$$ModuleName$$) Release(ctx context.Context) error {
	return nil
}
`

var tplModuleWire = `
package $$LowerModuleName$$

import (
	"$$ModuleImportPath$$/api"
	"$$ModuleImportPath$$/biz"
	"$$ModuleImportPath$$/dal"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new($$ModuleName$$), "*"),
)
`
