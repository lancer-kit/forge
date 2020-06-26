package project

import (
	"io/ioutil"
	"log"

	"github.com/go-ozzo/ozzo-validation"

	"gopkg.in/yaml.v2"
)

type ScaffoldTmplKey string

type tmplModuleOpts struct {
	Path string `yml:"path"`
}

func (s tmplModuleOpts) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Path, validation.Required),
	)
}

type TemplatesCfg struct {
	// Target defines the target directory info that is used by Scaffolder
	// to build all optional modules to its root path.
	Target tmplModuleOpts `yml:"target"`

	// Modules defines an optional service directories with the same directory
	// name for mapping the directories in base directory.
	Modules map[ScaffoldTmplKey]tmplModuleOpts `yml:"modules"`
}

func (s TemplatesCfg) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Target, validation.Required),
		validation.Field(&s.Modules, validation.Required),
	)
}

func ReadSchema(path string) TemplatesCfg {
	rawConfig, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read scaffold schema config file with path %s: %s", path, err)
	}

	config := new(TemplatesCfg)
	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		log.Fatalf("unable to scaffold schema config file with raw config %s: %s", rawConfig, err)
	}

	err = config.Validate()
	if err != nil {
		log.Fatalf("invalid scaffold schema config: %s", err)
	}

	return *config
}
