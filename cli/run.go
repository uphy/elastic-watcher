package cli

import (
	"os"

	"github.com/uphy/elastic-watcher/watcher"
	"github.com/urfave/cli"
)

func (c *CLI) run() cli.Command {
	return cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "print-config",
				Hidden: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			watchConf, err := watcher.LoadFile(ctx.Args().First())
			if err != nil {
				return err
			}
			if ctx.Bool("print-config") {
				watchConf.Save(os.Stdout)
			}
			w := watcher.NewWatch(c.globalConfig, watchConf)
			return w.Run()
		},
	}
}
