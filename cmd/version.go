package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Version(v string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show the version of the program",
		Action: func(c *cli.Context) error {
			fmt.Println(v)
			return nil
		},
	}
}
