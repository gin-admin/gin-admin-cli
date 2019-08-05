package cmd

import (
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
		},
		Action: func(c *cli.Context) error {
			return new.Exec(new.Config{
				Dir:       c.String("dir"),
				PkgName:   c.String("pkg"),
				UseMirror: c.Bool("mirror"),
			})
		},
	}
}
