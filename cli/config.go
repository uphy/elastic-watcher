package cli

import (
	"os"

	"github.com/urfave/cli"
)

func (c *CLI) config() cli.Command {
	return cli.Command{
		Name: "config",
		Action: func(ctx *cli.Context) error {
			c.globalConfig.Save(os.Stdout)
			return nil
		},
	}
}
