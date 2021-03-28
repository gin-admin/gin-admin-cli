package main

import (
	"log"
	"os"

	"github.com/gin-admin/gin-admin-cli/v4/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gin-admin-cli"
	app.Description = "gin-admin 辅助工具，提供创建项目、快速生成功能模块的功能"
	app.Version = "4.0.1"
	app.Commands = []cli.Command{
		cmd.NewCommand(),
		cmd.GenerateCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
