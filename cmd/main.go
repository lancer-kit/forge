package cmd

import (
	"strings"

	"github.com/lancer-kit/forge/configs"
	"github.com/urfave/cli"
)

func baseConfig(c *cli.Context) configs.BaseConfig {
	return configs.BaseConfig{
		Types:        strings.Split(c.String(typesFlag), ","),
		MergeSpecs:   c.Bool(mergeFlag),
		OutputPrefix: c.String(prefixFlag),
		OutputSuffix: c.String(suffixFlag),
		OutputDir:    c.String(dirFlag),
		OutputName:   c.String(nameFlag),
	}
}
