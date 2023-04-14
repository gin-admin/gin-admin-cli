package main

import (
	"embed"

	"github.com/gin-admin/gin-admin-cli/v10/internal/fs"
	"go.uber.org/zap"
)

//go:embed tpls
var f embed.FS

func main() {
	defer zap.S().Sync()

	// Set the embed.FS to the fs package
	fs.SetEFS(f)

	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

}
