package main

import (
	"os"

	"github.com/lancer-kit/forge/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "2.5"
	app.Name = "forge"
	app.Usage = "cli tool and generator from lancer-kit"
	app.Commands = cli.Commands{
		cmd.EnumCmd,
		cmd.ModelCmd,
		cmd.BindataCmd,
		cmd.ScaffoldCmd(),
	}

	_ = app.Run(os.Args)
}
