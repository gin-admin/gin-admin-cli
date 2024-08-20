package cmd

import (
	"github.com/gin-admin/gin-admin-cli/v10/internal/actions"
	"github.com/urfave/cli/v2"
)

func GenAntd() *cli.Command {
	return &cli.Command{
		Name:  "gen-antd",
		Usage: "Generate antd schemas to the specified module",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Usage:    "The project directory",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "The config file or directory",
			},
		},
		Action: func(c *cli.Context) error {
			gen := actions.GenAntd(actions.GenAntdConfig{
				Dir: c.String("dir"),
			})
			return gen.RunWithConfig(c.Context, c.String("config"))
		},
	}
}
