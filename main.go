package main

import (
	"os"

	"github.com/urfave/cli"
	"gitlab.inn4science.com/gophers/goplater/cmd"
)

func main() {
	app := cli.NewApp()
	app.Version = "2.4"
	app.Name = "goplater"
	app.Usage = "don't repeat yourself — generate from template"
	app.Commands = cli.Commands{
		cmd.EnumCmd,
		cmd.ModelCmd,
	}

	_ = app.Run(os.Args)
}
