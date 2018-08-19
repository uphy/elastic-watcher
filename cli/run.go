package cli

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/hashicorp/go-multierror"
	"github.com/uphy/elastic-watcher/pkg/watch"
)

func (c *CLI) run() cli.Command {
	return cli.Command{
		Name: "run",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:   "print-config",
				Hidden: true,
			},
			cli.BoolFlag{
				Name: "now",
			},
		},
		Action: func(ctx *cli.Context) error {
			path := ctx.Args().First()
			p, err := os.Stat(path)
			if os.IsNotExist(err) {
				return err
			}
			// collect watch config files
			ruleFiles := []string{}
			if p.IsDir() {
				d, err := ioutil.ReadDir(path)
				if err != nil {
					return err
				}
				for _, f := range d {
					ruleFiles = append(ruleFiles, filepath.Join(path, f.Name()))
				}
			} else {
				ruleFiles = append(ruleFiles, path)
			}
			// read watch config files
			watchConfigs := []*watch.WatchConfig{}
			for _, f := range ruleFiles {
				watchConf, err := watch.LoadFile(f)
				if err != nil {
					log.Printf("Failed to load config file %s: %v", f, err)
					continue
				}
				watchConfigs = append(watchConfigs, watchConf)
			}

			var errs error
			for _, conf := range watchConfigs {
				if ctx.Bool("print-config") {
					conf.Save(os.Stdout)
				}
				if ctx.Bool("now") {
					w := watch.NewWatch(c.globalConfig, conf)
					if err := w.Run(); err != nil {
						log.Printf("Failed to run watch: %v", err)
						errs = multierror.Append(errs, err)
					}
				}
			}
			if ctx.Bool("now") {
				return errs
			}

			wa := watch.New(c.globalConfig)
			for _, conf := range watchConfigs {
				if err := wa.AddWatch(conf); err != nil {
					log.Printf("Failed to add watch: %v", err)
				}
			}
			wa.Start()

			var sigc chan os.Signal
			sigc = make(chan os.Signal, 1)
			signal.Notify(sigc, os.Interrupt)
			<-sigc
			return nil
		},
	}
}
