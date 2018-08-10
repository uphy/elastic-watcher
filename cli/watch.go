package cli

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/uphy/elastic-watcher/watcher"
	"github.com/urfave/cli"
)

func (c *CLI) watch() cli.Command {
	return cli.Command{
		Name: "watch",
		Action: func(ctx *cli.Context) error {
			w := watcher.New(c.globalConfig)
			rulesDir := "rules"
			files, err := ioutil.ReadDir(rulesDir)
			if err != nil {
				return err
			}
			for _, f := range files {
				path := filepath.Join(rulesDir, f.Name())
				c, err := watcher.LoadFile(path)
				if err != nil {
					log.Printf("Failed to load config file %s: %v", path, err)
					continue
				}
				w.AddWatch(c)
			}
			w.Start()

			// wait for ctrl-c
			var signal_channel chan os.Signal
			signal_channel = make(chan os.Signal, 1)
			signal.Notify(signal_channel, os.Interrupt)
			<-signal_channel

			return nil
		},
	}
}
