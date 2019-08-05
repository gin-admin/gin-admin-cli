package main

import (
	"os"

	"github.com/LyricTian/gin-admin-cli/cmd"
	"github.com/LyricTian/logger"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gin-admin-cli"
	app.Description = "GinAdmin辅助工具"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		cmd.NewCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf(err.Error())
	}
}
