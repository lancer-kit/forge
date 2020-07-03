package configs

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type ScaffolderCfg struct {
	OutPath     string
	ProjectName string
	Schema      TemplatesCfg
	TmplModules ScaffoldTmplModules
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
