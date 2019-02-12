package main

import (
	"os"

	"github.com/urfave/cli"
	"gitlab.inn4science.com/gophers/forge/cmd"
)

func main() {
	app := cli.NewApp()
	app.Version = "2.4"
	app.Name = "forge"
	app.Usage = "cli tool and generator from lancer-kit"
	app.Commands = cli.Commands{
		cmd.EnumCmd,
		cmd.ModelCmd,
	}

	_ = app.Run(os.Args)
}
