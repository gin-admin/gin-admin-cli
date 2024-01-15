package cmd

import (
	"errors"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/actions"
	"github.com/gin-admin/gin-admin-cli/v10/internal/schema"
	"github.com/gin-admin/gin-admin-cli/v10/internal/tfs"
	"github.com/urfave/cli/v2"
)

// Generate returns the gen command.
func Generate() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"gen"},
		Usage:   "Generate structs to the specified module, support config file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Usage:    "The project directory to generate the struct",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "module",
				Aliases:  []string{"m"},
				Usage:    "The module to generate the struct from (like: RBAC)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "module-path",
				Usage: "The module path to generate the struct from (default: internal/mods)",
				Value: "internal/mods",
			},
			&cli.StringFlag{
				Name:  "wire-path",
				Usage: "The wire generate path to generate the struct from (default: internal/wirex)",
				Value: "internal/wirex",
			},
			&cli.StringFlag{
				Name:  "swag-path",
				Usage: "The swagger generate path to generate the struct from (default: internal/swagger)",
				Value: "internal/swagger",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "The config file or directory to generate the struct from (JSON/YAML)",
			},
			&cli.StringFlag{
				Name:  "structs",
				Usage: "The struct name to generate",
			},
			&cli.StringFlag{
				Name:  "structs-comment",
				Usage: "Specify the struct comment",
			},
			&cli.BoolFlag{
				Name:  "structs-router-prefix",
				Usage: "Use module name as router prefix",
			},
			&cli.StringFlag{
				Name:  "structs-output",
				Usage: "Specify the packages to generate the struct (default: schema,dal,biz,api)",
			},
			&cli.StringFlag{
				Name:  "tpl-path",
				Usage: "The template path to generate the struct from (default use tpls)",
			},
			&cli.StringFlag{
				Name:  "tpl-type",
				Usage: "The template type to generate the struct from (default: default)",
				Value: "default",
			},
			&cli.StringFlag{
				Name:  "fe-dir",
				Usage: "The frontend project directory to generate the UI",
			},
		},
		Action: func(c *cli.Context) error {
			if tplPath := c.String("tpl-path"); tplPath != "" {
				tfs.SetIns(tfs.NewOSFS(tplPath))
			}

			gen := actions.Generate(actions.GenerateConfig{
				Dir:         c.String("dir"),
				TplType:     c.String("tpl-type"),
				ModuleName:  c.String("module"),
				ModulePath:  c.String("module-path"),
				WirePath:    c.String("wire-path"),
				SwaggerPath: c.String("swag-path"),
				FEDir:       c.String("fe-dir"),
			})

			if c.String("config") != "" {
				return gen.RunWithConfig(c.Context, c.String("config"))
			} else if name := c.String("structs"); name != "" {
				var outputs []string
				if v := c.String("structs-output"); v != "" {
					outputs = strings.Split(v, ",")
				}
				return gen.RunWithStruct(c.Context, &schema.S{
					Name:             name,
					Comment:          c.String("structs-comment"),
					Outputs:          outputs,
					FillRouterPrefix: c.Bool("structs-router-prefix"),
				})
			} else {
				return errors.New("structs or config must be specified")
			}
		},
	}
}
