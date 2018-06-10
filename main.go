package main

import (
	"os"

	"github.com/sheb-gregor/goplater/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "2.0"
	app.Name = "goplater"
	app.Usage = "auto generate, don't repeat"
	app.Commands = cli.Commands{
		cmd.EnumCmd,
		cmd.ModelCmd,
	}

	_ = app.Run(os.Args)
}
