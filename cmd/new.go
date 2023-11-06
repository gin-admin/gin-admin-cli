package cmd

import (
	"context"

	"github.com/gin-admin/gin-admin-cli/v10/internal/actions"
	"github.com/urfave/cli/v2"
)

// New returns the new project command.
func New() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create a new project",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dir",
				Aliases:  []string{"d"},
				Usage:    "The directory to generate the project (default: current directory)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "The project name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "desc",
				Usage: "The project description",
			},
			&cli.StringFlag{
				Name:  "version",
				Usage: "The project version (default: 1.0.0)",
			},
			&cli.StringFlag{
				Name:  "pkg",
				Usage: "The project package name (default: project name)",
			},
			&cli.StringFlag{
				Name:  "git-url",
				Usage: "Use git repository to initialize the project (default: https://github.com/LyricTian/gin-admin.git)",
				Value: "https://github.com/LyricTian/gin-admin.git",
			},
			&cli.StringFlag{
				Name:  "git-branch",
				Usage: "Use git branch to initialize the project (default: main)",
				Value: "main",
			},
		},
		Action: func(c *cli.Context) error {
			n := actions.New(actions.NewConfig{
				Dir:         c.String("dir"),
				Name:        c.String("name"),
				Description: c.String("desc"),
				PkgName:     c.String("pkg"),
				Version:     c.String("version"),
				GitURL:      c.String("git-url"),
				GitBranch:   c.String("git-branch"),
			})
			return n.Run(context.Background())
		},
	}
}
