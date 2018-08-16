package main

import (
	"fmt"
	"os"

	"github.com/uphy/elastic-watcher/cli"
)

func main() {
	if err := cli.NewCLI().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
