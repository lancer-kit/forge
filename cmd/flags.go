package cmd

import (
	"github.com/urfave/cli"
)

const (
	typesFlag     = "type"
	tPathFlag     = "tmpl"
	mergeFlag     = "merge"
	suffixFlag    = "suffix"
	prefixFlag    = "prefix"
	dirFlag       = "dir"
	nameFlag      = "name"
	transformFlag = "transform"
	tprefixFlag   = "tprefix"
	tmplFlag      = "tmpl"
)

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
