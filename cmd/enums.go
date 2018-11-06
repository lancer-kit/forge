package cmd

import (
	"errors"
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/urfave/cli"
	"gitlab.inn4science.com/gophers/goplater/parser"
	"gitlab.inn4science.com/gophers/goplater/templates"
)

type enumsConfig struct {
	baseConfig
	transformRule string
	addTypePrefix bool
}

const (
	transformFlag = "transform"
	tprefixFlag   = "tprefix"
)

var EnumCmd = cli.Command{
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
	Action: genEnums,
}

func genEnums(c *cli.Context) error {
	config := enumsConfig{
		baseConfig: baseConfig{
			types:        strings.Split(c.String(typesFlag), ","),
			mergeSpecs:   c.Bool(mergeFlag),
			outputPrefix: c.String(prefixFlag),
			outputSuffix: c.String(suffixFlag),
		},

		transformRule: c.String(transformFlag),
		addTypePrefix: c.Bool(tprefixFlag),
	}

	// Only one directory at a time can be processed, and the default is ".".
	dir := "."
	if args := flag.Args(); len(args) == 1 {
		dir = args[0]
	} else if len(args) > 1 {
		return errors.New("only one directory at a time")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("unable to determine absolute filepath for requested path %s: %v",
			dir, err)
	}

	if len(config.types) == 1 {
		config.mergeSpecs = false
	}

	// need to remove already generated files for types
	// this is need for correct search of predefined by user
	// type vars and methods
	for _, typeName := range config.types {
		// Remove safe because we already check is path valid
		// and don't care about is present file - we need to remove it.
		os.Remove(config.getPath(typeName, dir))
	}

	if config.mergeSpecs {
		os.Remove(config.getPath(mergeTypeNames(config.types), dir))
	}

	pkg, err := parser.ParsePackage(dir)
	if err != nil {
		return fmt.Errorf("parsing package: %v", err)
	}

	var analysis = templates.Analysis{
		Command:     strings.Join(os.Args[1:], " "),
		PackageName: pkg.Name,
		Types:       make(map[string]templates.TypeSpec),
	}

	rule := templates.TransformRule(config.transformRule)

	// Run generate for each type.
	for _, typeName := range config.types {
		values, tmplsToExclude, err := pkg.ValuesOfType(typeName)
		if err != nil {
			return fmt.Errorf("finding values for type %v: %v", typeName, err)
		}
		analysis.Types[typeName] = templates.TypeSpec{
			TypeName:    typeName,
			Values:      rule.TransformValues(typeName, values, config.addTypePrefix),
			ExcludeList: tmplsToExclude,
		}
	}

	for name, src := range analysis.GenerateByTemplate(config.mergeSpecs) {
		if config.mergeSpecs {
			name = mergeTypeNames(config.types)
		}

		if err := ioutil.WriteFile(config.getPath(name, dir), src, 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}

		if config.mergeSpecs {
			return nil
		}
	}

	return nil
}

func mergeTypeNames(names []string) string {
	sort.Strings(names)
	single := strings.Join(names, "_")
	crc32InUint32 := crc32.ChecksumIEEE([]byte(single))
	fmt.Printf("NAMES :: %+v SINGLE %+v CRC32 :: %+v", names, single, strconv.FormatUint(uint64(crc32InUint32), 16))
	return strconv.FormatUint(uint64(crc32InUint32), 16)
}
