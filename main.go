package main

import (
	"log"
	"os"

	"github.com/gin-admin/gin-admin-cli/v6/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gin-admin-cli"
	app.Description = "gin-admin v9 generate tools (create project and generate modules)"
	app.Version = "6.0.2"
	app.Commands = []cli.Command{
		cmd.NewCommand(),
		cmd.GenerateCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
