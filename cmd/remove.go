package cmd

import (
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/actions"
	"github.com/urfave/cli/v2"
)

// Remove returns the remove command.
func Remove() *cli.Command {
	return &cli.Command{
		Name:    "remove",
		Aliases: []string{"rm"},
		Usage:   "Remove structs from the module",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Usage:    "The directory to remove the struct from",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "module",
				Aliases:  []string{"m"},
				Usage:    "The module to remove the struct from",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "module-path",
				Usage: "The module path to remove the struct from (default: internal/mods)",
				Value: "internal/mods",
			},
			&cli.StringFlag{
				Name:    "structs",
				Aliases: []string{"s"},
				Usage:   "The struct to remove (multiple structs can be separated by a comma)",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "The config file to generate the struct from (JSON/YAML)",
			},
			&cli.StringFlag{
				Name:  "wire-path",
				Usage: "The wire generate path to remove the struct from (default: internal/library/wirex)",
				Value: "internal/wirex",
			},
			&cli.StringFlag{
				Name:  "swag-path",
				Usage: "The swagger generate path to remove the struct from (default: internal/swagger)",
				Value: "internal/swagger",
			},
		},
		Action: func(c *cli.Context) error {
			rm := actions.NewRemove(&actions.RemoveConfig{
				Dir:         c.String("dir"),
				ModuleName:  c.String("module"),
				ModulePath:  c.String("module-path"),
				WirePath:    c.String("wire-path"),
				SwaggerPath: c.String("swag-path"),
			})

			if c.String("config") != "" {
				return rm.RunWithConfig(c.Context, c.String("config"))
			} else if c.String("structs") != "" {
				return rm.Run(c.Context, strings.Split(c.String("structs"), ","))
			}
			return nil
		},
	}
}
