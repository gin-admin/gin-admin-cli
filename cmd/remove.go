package cmd

import "github.com/urfave/cli/v2"

func Remove() *cli.Command {
	return &cli.Command{
		Name:  "remove",
		Usage: "Remove multiple structs from the module",
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
				Name:     "structs",
				Aliases:  []string{"s"},
				Usage:    "The struct to remove (multiple structs can be separated by a comma)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {

			return nil
		},
	}
}
