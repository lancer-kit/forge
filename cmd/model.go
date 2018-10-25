package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"gitlab.inn4science.com/gophers/goplater/parser"
	"gitlab.inn4science.com/gophers/goplater/templates"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// need to remove already generated files for types
	// this is need for correct search of predefined by user
	// type vars and methods
	for _, typeName := range config.types {
		// Remove safe because we already check is path valid
		// and don't care about is present file - we need to remove it.
		os.Remove(config.getPath(typeName, dir))
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

	tmpl, err := templates.OpenTemplate(config.tPath)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Run generate for each type.
	for _, typeName := range config.types {
		spec, err := pkg.FindStructureSpec(typeName)
		if err != nil {
			return fmt.Errorf("finding values for type %v: %v", typeName, err)
		}

		model := templates.FigureOut(spec)
		model.Package = pkg.Name

		newRawFile, err := model.Exec(tmpl)
		if err != nil {
			log.Fatalf("exec template for type %v failed: %v", typeName, err)
		}

		if err := ioutil.WriteFile(config.getPath(typeName, dir), []byte(newRawFile), 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}
	}

	return nil
}
