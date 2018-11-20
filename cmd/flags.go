package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/urfave/cli"
)

const (
	typesFlag  = "type"
	mergeFlag  = "merge"
	suffixFlag = "suffix"
	prefixFlag = "prefix"
)

type BaseConfig struct {
	types        []string
	mergeSpecs   bool
	outputSuffix string
	outputPrefix string
}

func (BaseConfig) FromContext(c *cli.Context) BaseConfig {
	return BaseConfig{
		types:        strings.Split(c.String(typesFlag), ","),
		mergeSpecs:   c.Bool(mergeFlag),
		outputPrefix: c.String(prefixFlag),
		outputSuffix: c.String(suffixFlag),
	}
}

func (config *BaseConfig) Validate() error {
	if len(config.types) == 0 {
		return fmt.Errorf("%s: should not be empty", typesFlag)
	}
	if config.outputPrefix == "" && config.outputSuffix == "" {
		return fmt.Errorf("%s or %s: should be passed", suffixFlag, prefixFlag)
	}

	return nil
}

func (config BaseConfig) GetPath(name, dir string) string {
	var splittedName string

	for i, r := range name {
		if i == 0 {
			splittedName += string(r)
			continue
		}
		if unicode.IsUpper(r) {
			splittedName += "_" + string(r)
			continue
		}
		splittedName += string(r)
	}

	output := strings.ToLower(config.outputPrefix + splittedName + config.outputSuffix + ".go")

	return filepath.Join(dir, output)
}

var baseFlags = []cli.Flag{
	cli.StringFlag{
		Name:  typesFlag,
		Usage: "list of type names; required;",
	},
	cli.StringFlag{
		Name:  prefixFlag,
		Usage: "prefix to be added to the output file;",
		Value: "enums_",
	},

	cli.StringFlag{
		Name:  suffixFlag,
		Usage: "suffix to be added to the output file;",
		Value: "",
	},

	cli.BoolFlag{
		Name:  mergeFlag,
		Usage: "merge all output into one file;",
	},
}
