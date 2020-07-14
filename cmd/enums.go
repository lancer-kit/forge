package cmd

import (
	"github.com/urfave/cli"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/generate"
	"github.com/lancer-kit/forge/templates"
)

func EnumCmd() cli.Command {
	return cli.Command{
		Name:  "enum",
		Usage: "generate var and methods for the iota-enums",
		Flags: append(baseFlags,

			cli.StringFlag{
				Name:  transformFlag,
				Usage: "way to convert constants to a string;",
				Value: "none",
			},

			cli.BoolFlag{
				Name:  tprefixFlag,
				Usage: "keep typename prefix in string values or not;",
			},
		),
		Action: enumsAction,
	}
}

func enumsAction(c *cli.Context) error {
	config := enumsConfig(c)
	err := config.Validate()
	if err != nil {
		return cli.NewExitError("ERROR: "+err.Error(), 1)
	}
	err = generate.Enums(config)
	if err != nil {
		return cli.NewExitError("ERROR: "+err.Error(), 1)
	}
	return nil
}

func enumsConfig(c *cli.Context) configs.EnumsConfig {
	return configs.EnumsConfig{
		BaseConfig:    baseConfig(c),
		TransformRule: templates.TransformRule(c.String(transformFlag)),
		AddTypePrefix: c.Bool(tprefixFlag),
	}
}
