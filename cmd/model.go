package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sheb-gregor/goplater/parser"
	"github.com/sheb-gregor/goplater/templates"
	"github.com/urfave/cli"
)

const tPath = "tmpl"

var ModelCmd = cli.Command{
	Name:  "model",
	Usage: "generate code for structure by template",
	Flags: append(baseFlags,
		cli.StringFlag{
			Name:  tPath,
			Usage: "path to the templates; required;",
		},
	),
	Action: genModelAction,
}

type stConfig struct {
	baseConfig
	tPath string
}

func genModelAction(c *cli.Context) error {
	config := stConfig{
		baseConfig: baseConfig{
			types:        strings.Split(c.String(typesFlag), ","),
			mergeSpecs:   c.Bool(mergeFlag),
			outputPrefix: c.String(prefixFlag),
			outputSuffix: c.String(suffixFlag),
		},

		tPath: c.String(tPath),
	}
	err := genModel(config)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func genModel(config stConfig) error {

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

	file, err := os.Open(config.tPath)
	if err != nil {
		return fmt.Errorf("unable to open template file [%s]: %v",
			config.tPath, err)
	}
	defer file.Close()
	rawTemplate := make([]byte, 0)
	_, err = file.Read(rawTemplate)
	if err != nil {
		return fmt.Errorf("unable to read template file [%s]: %v",
			config.tPath, err)
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

	_ = templates.Analysis{
		Command:     strings.Join(os.Args[1:], " "),
		PackageName: pkg.Name,
		Types:       make(map[string]templates.TypeSpec),
	}

	// Run generate for each type.
	for _, typeName := range config.types {
		res, err := pkg.FindStructureSpec(typeName)
		if err != nil {
			return fmt.Errorf("finding values for type %v: %v", typeName, err)
		}

		tmpl := templates.StructSpec{
			Name: typeName,
		}
		for _, value := range res.Fields {
			tmpl.Fields = append(tmpl.Fields, templates.Field{
				Name: value,
				Type: res.FTypes[value],
				Tags: parseStructTags(res.Tags[value]),
			})
		}

		fmt.Printf("%+v\n", res)
		fmt.Printf("%+v\n", tmpl)
	}

	return nil
}

func parseStructTags(tag string) map[string]string {
	tag = strings.Trim(tag, "`")
	tags := map[string]string{}
	for _, fullTag := range strings.Split(tag, " ") {
		tag := strings.Split(fullTag, ":")
		tags[tag[0]] = strings.Trim(strings.Split(tag[1], ",")[0], `"`)
	}
	return tags
}
