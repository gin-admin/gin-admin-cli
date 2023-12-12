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
				Name:  "app-name",
				Usage: "The application name (default: project name)",
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
			&cli.StringFlag{
				Name:  "fe-dir",
				Usage: "The frontend directory to generate the project (if empty, the frontend project will not be generated)",
			},
			&cli.StringFlag{
				Name:  "fe-name",
				Usage: "The frontend project name (default: frontend)",
				Value: "frontend",
			},
			&cli.StringFlag{
				Name:  "fe-git-url",
				Usage: "Use git repository to initialize the frontend project (default: https://github.com/gin-admin/gin-admin-frontend.git)",
				Value: "https://github.com/gin-admin/gin-admin-frontend.git",
			},
			&cli.StringFlag{
				Name:  "fe-git-branch",
				Usage: "Use git branch to initialize the frontend project (default: main)",
				Value: "main",
			},
		},
		Action: func(c *cli.Context) error {
			n := actions.New(actions.NewConfig{
				Dir:         c.String("dir"),
				Name:        c.String("name"),
				AppName:     c.String("app-name"),
				Description: c.String("desc"),
				PkgName:     c.String("pkg"),
				Version:     c.String("version"),
				GitURL:      c.String("git-url"),
				GitBranch:   c.String("git-branch"),
				FeDir:       c.String("fe-dir"),
				FeName:      c.String("fe-name"),
				FeGitURL:    c.String("fe-git-url"),
				FeGitBranch: c.String("fe-git-branch"),
			})
			return n.Run(context.Background())
		},
	}
}
