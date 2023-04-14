package main

import (
	"embed"
	"os"

	"github.com/gin-admin/gin-admin-cli/v10/cmd"
	"github.com/gin-admin/gin-admin-cli/v10/internal/tfs"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

//go:embed tpls
var f embed.FS

var VERSION = "v10.0.0-beta"

func main() {
	defer zap.S().Sync()

	// Set the embed.FS to the fs package
	tfs.SetEFS(f)

	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	app := cli.NewApp()
	app.Name = "gin-admin-cli"
	app.Version = VERSION
	app.Usage = "gin-admin-cli is a command line tool for gin-admin."
	app.Authors = append(app.Authors, &cli.Author{
		Name:  "LyricTian",
		Email: "tiannianshou@gmail.com",
	})
	app.Commands = []*cli.Command{
		cmd.Version(VERSION),
		cmd.Generate(),
		cmd.Remove(),
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
