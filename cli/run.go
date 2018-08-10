package cli

import (
	"github.com/uphy/elastic-watcher/watcher"
	"github.com/urfave/cli"
)

func (c *CLI) run() cli.Command {
	return cli.Command{
		Name: "run",
		Action: func(ctx *cli.Context) error {
			watchConf, err := watcher.LoadFile(ctx.Args().First())
			if err != nil {
				return err
			}
			//watchConf.Save(os.Stdout)
			w := watcher.NewWatch(c.globalConfig, watchConf)
			return w.Run()
		},
	}
}
