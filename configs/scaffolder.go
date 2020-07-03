package configs

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

type ScaffolderCfg struct {
	OutPath     string
	ProjectName string
	Schema      TmplSchemaCfg
	TmplModules ScaffoldTmplModules

	// TmplName defines the key name of the template that is predefined in
	// schema.yml file. In case if the TmplName is empty the base template
	// is grabbed from schema.yml file
	TmplName string
}

func (cfg ScaffolderCfg) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.OutPath, validation.Required),
		validation.Field(&cfg.Schema, validation.Required),
		validation.Field(&cfg.ProjectName, validation.Required),
		validation.Field(&cfg.TmplModules, validation.Required),
	)
}

type ScaffoldTmplKey string

type tmplModuleOpts struct {
	Path string `yml:"path"`
}

func (s tmplModuleOpts) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Path, validation.Required),
	)
}

type TmplSchemaCfg struct {
	// Base contains the key name that is defined in Specs cfg
	Base string `yml:"base"`

	// Specs contains the templates specification where
	// the key is the template name and SpecCfg defines the
	// template behaviour as a root directory with base template body
	// and submodules directory info for optional scaffold
	Specs map[string]SpecCfg `yml:"specs"`
}

func (cfg TmplSchemaCfg) Validate() error {
	err := validation.ValidateStruct(&cfg,
		validation.Field(&cfg.Base, validation.Required),
		validation.Field(&cfg.Specs, validation.Required),
	)
	if err != nil {
		return err
	}

	_, ok := cfg.Specs[cfg.Base]
	if !ok {
		return errors.New("no base specification predefined in template schema")
	}
	return nil
}

type SpecCfg struct {
	Path string `yml:"path"`
	// Target defines the target directory info that is used by Scaffolder
	// to build all optional modules to its root path.
	Target tmplModuleOpts `yml:"target"`

	// Modules defines an optional service directories with the same directory
	// name for mapping the directories in base directory.
	Modules map[ScaffoldTmplKey]tmplModuleOpts `yml:"modules"`
}

func (cfg SpecCfg) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.Path, validation.Required),
		validation.Field(&cfg.Target, validation.Required),
	)
}

type ScaffoldTmplModules map[interface{}]interface{}

const (
	ScaffoldProjectNameKey = "project_name"

	// Module key names for optional scaffolding
	// !Don`t change the name of the keys as they must be the same in
	// - all templates .go.tpl
	// - scaffold schema schema.yml
	// - scaffold data (passes to render the templates)
	ModuleKeyAPI          ScaffoldTmplKey = "api"
	ModuleKeyDB           ScaffoldTmplKey = "db"
	ModuleKeySimpleWorker ScaffoldTmplKey = "simple_worker"
)
