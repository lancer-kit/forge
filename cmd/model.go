package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const tPath = "tmpl"

var ModelCmd = cli.Command{
	Name:  "model",
	Usage: "generate code for structure by template",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  tPath,
			Usage: "path to the templates; required;",
		},
		cli.StringFlag{
			Name:  typesFlag,
			Usage: "list of type names; required;",
		},
		cli.StringFlag{
			Name:  prefixFlag,
			Usage: "prefix to be added to the output file;",
			Value: "",
		},

		cli.StringFlag{
			Name:  suffixFlag,
			Usage: "suffix to be added to the output file;",
			Value: "",
		},
	},
	Action: genModelAction,
}

func (config *ModelConfig) Validate() error {
	if err := config.BaseConfig.Validate(); err != nil {
		return err
	}
	if config.tPath == "" {
		return fmt.Errorf("%s: must be specified", tPath)
	}

	_, err := os.Stat(config.tPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: file is not exist", tPath)
	}
	if err != nil {
		return fmt.Errorf("%s: %s", tPath, err.Error())
	}
	return nil
}

func genModelAction(c *cli.Context) error {
	config := ModelConfig{}.FromContext(c)
	if err := config.Validate(); err != nil {
		return cli.NewExitError("[ERROR] "+err.Error(), 1)
	}

	err := genModel(config)
	if err != nil {
		return cli.NewExitError("[ERROR] "+err.Error(), 1)
	}

	return nil
}
