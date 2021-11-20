package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-admin/gin-admin-cli/v6/cmd/generate"
	"github.com/gin-admin/gin-admin-cli/v6/cmd/new"
	"github.com/urfave/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir, d",
				Usage: "Project directory (default: GOPATH/src/package_name)",
			},
			&cli.StringFlag{
				Name:  "pkg, p",
				Usage: "Package name",
			},
			&cli.StringFlag{
				Name:  "branch, b",
				Usage: "Git branch",
				Value: "master",
			},
			&cli.BoolFlag{
				Name:  "mirror, m",
				Usage: "Use gitee (gitee.com)",
			},
			&cli.BoolFlag{
				Name:  "tpl",
				Usage: "Use gin-admin-tpl",
			},
			&cli.BoolFlag{
				Name:  "web, w",
				Usage: "Include gin-admin-react",
			},
		},
		Action: func(c *cli.Context) error {
			cfg := new.Config{
				Dir:        c.String("dir"),
				PkgName:    c.String("pkg"),
				UseMirror:  c.Bool("mirror"),
				UseTpl:     c.Bool("tpl"),
				Branch:     c.String("branch"),
				IncludeWeb: c.Bool("web"),
			}

			if cfg.PkgName == "" {
				return errors.New("请指定包名")
			}

			if cfg.Dir == "" {
				vpath := os.Getenv("GOPATH")
				if vpath == "" {
					return errors.New("please specify project directory")
				}
				cfg.Dir = filepath.Join(vpath, "src", cfg.PkgName)
			}

			cfg.AppName = filepath.Base(cfg.PkgName)

			return new.Exec(cfg)
		},
	}
}

func GenerateCommand() cli.Command {
	return cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir, d",
				Usage: "Project directory(default: GOPATH)",
			},
			&cli.StringFlag{
				Name:     "pkg, p",
				Usage:    "Package name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "name, n",
				Usage: "Struct name",
			},
			&cli.StringFlag{
				Name:  "comment, c",
				Usage: "Struct comment",
			},
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "Template file (.yaml)",
			},
			&cli.StringFlag{
				Name:  "module, m",
				Usage: "Specify generate modules (schema,dao,service,api,router）",
			},
			&cli.BoolFlag{
				Name:  "include_status",
				Usage: "whether include status field",
			},
			&cli.BoolFlag{
				Name:  "include_creator",
				Usage: "whether include created_by field",
			},
		},
		Action: func(c *cli.Context) error {
			cfg := generate.Config{
				Dir:           c.String("dir"),
				PkgName:       c.String("pkg"),
				Name:          c.String("name"),
				Comment:       c.String("comment"),
				File:          c.String("file"),
				Modules:       c.String("module"),
				ExcludeStatus: !c.Bool("include_status"),
				ExcludeCreate: !c.Bool("include_creator"),
			}

			if cfg.Dir == "" {
				vpath := os.Getenv("GOPATH")
				if vpath == "" {
					return errors.New("please specify project directory")
				}
				cfg.Dir = filepath.Join(vpath, "src", cfg.PkgName)
			}

			if cfg.PkgName == "" {
				fmt.Println("Package name not be empty")
				return nil
			} else if cfg.Name == "" && cfg.File == "" {
				fmt.Println("Please specify struct name or template file")
				return nil
			}

			return generate.Exec(cfg)
		},
	}
}
