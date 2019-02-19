package cmd

import (
	"strings"

	"github.com/urfave/cli"
	"gitlab.inn4science.com/gophers/forge/configs"
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
