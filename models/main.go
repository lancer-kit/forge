package models

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/sheb-gregor/goplater/parser"
	"github.com/sheb-gregor/goplater/templates"
)

// todo: refactor
// 1. Analyze model(s)
// 2. Gen by template
// 3. Write file
func GenModel(config ModelConfig) error {
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
			log.Printf("[WARN] definition of the type %s isn't found, skip it. \n", typeName)
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
