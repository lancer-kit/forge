package cmd

import (
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

const (
	typesFlag  = "type"
	mergeFlag  = "merge"
	suffixFlag = "suffix"
	prefixFlag = "prefix"
)

type baseConfig struct {
	types        []string
	mergeSpecs   bool
	outputSuffix string
	outputPrefix string
}

func (config baseConfig) getPath(name, dir string) string {
	output := strings.ToLower(config.outputPrefix + name + config.outputSuffix + ".go")
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
