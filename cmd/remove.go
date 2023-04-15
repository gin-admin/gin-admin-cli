package cmd

import (
	"github.com/gin-admin/gin-admin-cli/v10/internal/actions"
	"github.com/urfave/cli/v2"
)

// Remove returns the remove command.
func Remove() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "Remove structs from the module",
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
				Name:     "structs",
				Aliases:  []string{"s"},
				Usage:    "The struct to remove (multiple structs can be separated by a comma)",
				Required: true,
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
			return actions.NewRemove(&actions.RemoveConfig{
				Dir:         c.String("dir"),
				ModuleName:  c.String("module"),
				ModulePath:  c.String("module-path"),
				WirePath:    c.String("wire-path"),
				SwaggerPath: c.String("swag-path"),
			}).Run(c.Context, c.String("structs"))
		},
	}
}
