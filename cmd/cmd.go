package cmd

import (
	"fmt"

	"github.com/LyricTian/gin-admin-cli/cmd/generate"
	"github.com/LyricTian/gin-admin-cli/cmd/new"
	"github.com/urfave/cli"
)

// NewCommand 创建项目命令
func NewCommand() cli.Command {
	return cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir, d",
				Usage: "项目生成目录",
			},
			&cli.StringFlag{
				Name:  "pkg, p",
				Usage: "项目包名",
			},
			&cli.BoolFlag{
				Name:  "mirror, m",
				Usage: "使用国内镜像(gitee.com)",
			},
			&cli.BoolFlag{
				Name:  "web, w",
				Usage: "包含web项目",
			},
		},
		Action: func(c *cli.Context) error {
			cfg := new.Config{
				Dir:        c.String("dir"),
				PkgName:    c.String("pkg"),
				UseMirror:  c.Bool("mirror"),
				IncludeWeb: c.Bool("web"),
			}
			return new.Exec(cfg)
		},
	}
}

// GenerateCommand 生成项目模块命令
func GenerateCommand() cli.Command {
	return cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir, d",
				Usage:    "项目生成目录",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pkg, p",
				Usage:    "项目包名",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "ctl",
				Usage: "控制器swagger模板(支持default(基于github.com/swaggo/swag)和tb(基于github.com/teambition/swaggo))",
				Value: "default",
			},
			&cli.StringFlag{
				Name:  "router",
				Usage: "路由模块(routers/api/api.go)",
				Value: "api",
			},
			&cli.StringFlag{
				Name:  "name, n",
				Usage: "业务模块名称(结构体名称)",
			},
			&cli.StringFlag{
				Name:  "comment, c",
				Usage: "业务模块注释(结构体注释)",
			},
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "指定模板文件(.json，模板配置可参考说明)",
			},
			&cli.StringFlag{
				Name:  "module, m",
				Usage: "指定生成模块（以逗号分隔，支持：all,schema,entity,model,bll,router）",
			},
		},
		Action: func(c *cli.Context) error {
			cfg := generate.Config{
				Dir:        c.String("dir"),
				PkgName:    c.String("pkg"),
				CtlTpl:     c.String("ctl"),
				RouterName: c.String("router"),
				Name:       c.String("name"),
				Comment:    c.String("comment"),
				File:       c.String("file"),
				Modules:    c.String("module"),
			}

			if cfg.Dir == "" {
				fmt.Println("请指定项目目录")
				return nil
			} else if cfg.PkgName == "" {
				fmt.Println("请指定包名")
				return nil
			} else if cfg.Name == "" && cfg.File == "" {
				fmt.Println("请指定模块名称或模板文件")
				return nil
			} else if cfg.Name != "" && cfg.Comment == "" {
				fmt.Println("请指定模块说明")
				return nil
			}

			return generate.Exec(cfg)
		},
	}
}
