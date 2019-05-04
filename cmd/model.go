package cmd

import (
	"github.com/urfave/cli"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/generate"
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
	Action: modelAction,
}

func modelAction(c *cli.Context) error {
	config := modelConfig(c)
	if err := config.Validate(); err != nil {
		return cli.NewExitError("[ERROR] "+err.Error(), 1)
	}

	err := generate.Model(config)
	if err != nil {
		return cli.NewExitError("[ERROR] "+err.Error(), 1)
	}

	return nil
}

func modelConfig(c *cli.Context) configs.ModelConfig {
	return configs.ModelConfig{
		BaseConfig: baseConfig(c),
		TPath:      c.String(tmplFlag),
	}
}
