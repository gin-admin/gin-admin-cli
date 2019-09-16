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
				Name:  "name, n",
				Usage: "模块名称(结构体名称)",
			},
			&cli.StringFlag{
				Name:  "comment, c",
				Usage: "模块注释",
			},
			&cli.StringFlag{
				Name:  "file, f",
				Usage: "指定模板文件(.json，模板配置可参考说明)",
			},
		},
		Action: func(c *cli.Context) error {
			cfg := generate.Config{
				Dir:     c.String("dir"),
				PkgName: c.String("pkg"),
				Name:    c.String("name"),
				Comment: c.String("comment"),
				File:    c.String("file"),
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
