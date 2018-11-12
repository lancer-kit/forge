package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"gitlab.inn4science.com/gophers/goplater/parser"
	"gitlab.inn4science.com/gophers/goplater/templates"
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

type ModelConfig struct {
	BaseConfig
	tPath string
}

func (ModelConfig) FromContext(c *cli.Context) ModelConfig {
	return ModelConfig{
		BaseConfig: BaseConfig{}.FromContext(c),
		tPath:      c.String(tPath),
	}
}

func (config *ModelConfig) Validate() error {
	if err := config.BaseConfig.Validate(); err != nil {
		return err
	}
	if config.tPath == "" {
		return fmt.Errorf("%s: must be specified", tPath)
	}

	_, err := os.Stat(config.tPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: file is not exist", tPath)
	}
	if err != nil {
		return fmt.Errorf("%s: %s", tPath, err.Error())
	}
	return nil
}

func genModelAction(c *cli.Context) error {
	config := ModelConfig{}.FromContext(c)
	if err := config.Validate(); err != nil {
		return cli.NewExitError("ERROR: "+err.Error(), 1)
	}

	err := genModel(config)
	if err != nil {
		return cli.NewExitError("ERROR: "+err.Error(), 1)
	}

	return nil
}

func genModel(config ModelConfig) error {
	// Only one directory at a time can be processed, and the default is ".".
	dir := "."

	if args := flag.Args(); len(args) >= 1 {
		dir = args[0]
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
		os.Remove(config.GetPath(typeName, dir))
	}

	pkg, err := parser.ParsePackage(dir)
	if err != nil {
		return fmt.Errorf("parsing package: %v", err)
	}

	tmpl, err := templates.OpenTemplate(config.tPath)
	if err != nil {
		return fmt.Errorf("unable to open template: %s", err.Error())
	}

	// Run generate for each type.
	for _, typeName := range config.types {
		spec, err := pkg.FindStructureSpec(typeName)
		if err != nil {
			return fmt.Errorf("finding values for type %v: %s", typeName, err.Error())
		}
		if spec == nil {
			fmt.Printf("WARN: definition of the type %s isn't found, skip it. \n", typeName)
			continue
		}

		model, err := templates.FigureOut(spec)
		if err != nil {
			return fmt.Errorf("FigureOut for type %v is failed: %s", typeName, err.Error())
		}

		model.Package = pkg.Name
		newRawFile, err := model.Exec(tmpl)
		if err != nil {
			return fmt.Errorf("exec template for type %v failed: %v", typeName, err)
		}

		if err := ioutil.WriteFile(config.GetPath(typeName, dir), []byte(newRawFile), 0644); err != nil {
			return fmt.Errorf("writing output failed: %s", err)
		}
	}

	return nil
}
