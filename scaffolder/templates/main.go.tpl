package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"{{.project_name}}/cmd"
	"{{.project_name}}/config"
	"{{.project_name}}/info"
)

func main() {
	app := cli.NewApp()
	app.Usage = "A " + config.ServiceName + " service"
	app.Version = info.App.Version
	app.Flags = cmd.GetFlags()
	app.Commands = cmd.GetCommands()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
