package cmd

import (
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

type baseConfig struct {
	types        []string
	mergeSpecs   bool
	outputSuffix string
	outputPrefix string
}

func (config baseConfig) getPath(name, dir string) string {
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
