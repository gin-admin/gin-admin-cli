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

var VERSION = "v10.3.3"

func main() {
	defer func() {
		_ = zap.S().Sync()
	}()

	// Set the embed.FS to the fs package
	tfs.SetEFS(f)

	logger, err := zap.NewDevelopmentConfig().Build(zap.WithCaller(false))
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	app := cli.NewApp()
	app.Name = "gin-admin-cli"
	app.Version = VERSION
	app.Usage = "A command line tool for [gin-admin](https://github.com/LyricTian/gin-admin)."
	app.Authors = append(app.Authors, &cli.Author{
		Name:  "LyricTian",
		Email: "tiannianshou@gmail.com",
	})
	app.Commands = []*cli.Command{
		cmd.Version(VERSION),
		cmd.New(),
		cmd.Generate(),
		cmd.Remove(),
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
