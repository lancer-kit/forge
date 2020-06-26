package project

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
