package main

import (
	"log"
	"os"

	"github.com/gin-admin/gin-admin-cli/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gin-admin-cli"
	app.Description = "GinAdmin辅助工具，提供创建项目、快速生成功能模块的功能"
	app.Version = "3.1.0"
	app.Commands = []cli.Command{
		cmd.NewCommand(),
		cmd.GenerateCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
