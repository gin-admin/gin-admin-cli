package parser

var tplModuleMain = `
package $$LowerModuleName$$

import (
	"context"

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
	if err := a.AutoMigrate(ctx); err != nil {
		return err
	}
	return nil
}

func (a *$$ModuleName$$) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
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

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new($$ModuleName$$), "*"),
)
`
