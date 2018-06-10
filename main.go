package main

import (
	"os"

	"github.com/urfave/cli"
)

const (
	typesFlag     = "type"
	transformFlag = "transform"
	tprefixFlag   = "tprefix"
	mergeFlag     = "merge"
	suffixFlag    = "suffix"
	prefixFlag    = "prefix"
)

func main() {
	app := cli.NewApp()
	app.Version = "2.0"
	app.Name = "goplater"
	app.Usage = "auto generate, don't repeat"
	app.Commands = cli.Commands{
		enumCmd,
	}

	_ = app.Run(os.Args)
}
