package srv

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"

	forge "github.com/lancer-kit/forge/.forge"
	"github.com/lancer-kit/forge/configs"
)

const (
	TmplTypeDefault string = "default"
	TmplTypeForge   string = "forge"

	ProjectGenGoPathType    string = "Go Path"
	ProjectGenGoModulesType string = "Go Modules"
)

// Cfg represents srv cfg for scaffolding the project
type Cfg struct {
	// Template type represents the template provider.
	// Can be default - is embedded into package and .forge templates
	TemplateType string

	// ProjectGenType represents project generation type. Can be
	// Go Path or Go Modules
	ProjectGenType string

	GitOrigin string //*

	// GoPathDomainName represents project domain name path in GOPATH dir
	GoPathDomainName string // *

	// OutDir represents the output dir in case of Go Modules project scaffold
	OutDir string // *

	GoModulesProjectName string // *

	// ForgeTmplKeyName defines the template key name defined in .forge
	// directory in schema.yml file
	ForgeTmplKeyName string
}

func AskSurvey() (*Cfg, error) {
	var (
		surveyCfg = new(Cfg)
		err       error
	)

	surveyCfg.ProjectGenType, err = askPromptSelect("Choose project scaffold type:", TmplTypeDefault,
		[]string{TmplTypeDefault, TmplTypeForge})
	if err != nil {
		return nil, fmt.Errorf("failed to get the survey answer: %s", err)
	}

	switch surveyCfg.ProjectGenType {
	case TmplTypeDefault:
		err = surveyCfg.ask()
		if err != nil {
			return nil, err
		}
		return surveyCfg, nil

	case TmplTypeForge:
		// extract all forge template key names
		tmpls, err := getForgeTmplKeyNames()
		if err != nil {
			return nil, err
		}
		surveyCfg.ForgeTmplKeyName, err = askPromptSelect("Choose forge template name:", "", tmpls)
		if err != nil {
			return nil, fmt.Errorf("failed to get the forge template name: %s", err)
		}
		err = surveyCfg.ask()
		if err != nil {
			return nil, err
		}

		return surveyCfg, nil
	}
	return surveyCfg, err
}

func (c *Cfg) ask() error {
	var (
		genTypeSurvey = new(genTypeSurvey)
		err           error
	)

	genTypeSurvey, err = askProjectGenTypeQuestion()
	if err != nil {
		return fmt.Errorf("failed to get the survey answer: %s", err)
	}
	if genTypeSurvey.WithGitOrigin == "yes" {
		c.GitOrigin, err = askPromptInput("Enter Git origin address:", "")
		if err != nil {
			return fmt.Errorf("failed to get the git origin address")
		}
	}

	switch genTypeSurvey.ProjectGenType {
	case ProjectGenGoPathType:
		c.GoPathDomainName, err = askPromptInput("Enter Project domain name` "+
			"(ex. gitlab.com/team/project) to be created in GOPATH", "forge")
		if err != nil {
			return fmt.Errorf("failed to get the survey answer: %s", err)

		}
		return nil

	case ProjectGenGoModulesType:
		var genModulesSurvey = new(genGoModulesSurvey)
		genModulesSurvey, err = askGenGoModulesSurvey()
		if err != nil {
			return fmt.Errorf("failed to get the survey answer: %s", err)
		}

		c.OutDir = genModulesSurvey.OutDirPath
		c.GoModulesProjectName = genModulesSurvey.ProjectName

		return nil
	default:
		return fmt.Errorf("wron project generation type")
	}
}

func getForgeTmplKeyNames() ([]string, error) {
	asset, err := forge.Asset(configs.ForgeSchemaAssetName)
	if err != nil {
		return nil, fmt.Errorf("failed to load the %s asset: %s", configs.ForgeSchemaAssetName, err)
	}

	forgeSchema := map[string]configs.ForgeTmpl{}
	err = yaml.Unmarshal(asset, &forgeSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the forge schema: %s", err)
	}

	keys := make([]string, 0, len(forgeSchema))

	for k := range forgeSchema {
		keys = append(keys, k)
	}
	return keys, nil
}

func (c *Cfg) GetForgeTmpl() (*configs.ForgeTmpl, error) {
	if c.ForgeTmplKeyName != "" {
		// Get Forge schema to config
		asset, err := forge.Asset(configs.ForgeSchemaAssetName)
		if err != nil {
			return nil, fmt.Errorf("failed to load the %s asset: %s", configs.ForgeSchemaAssetName, err)
		}

		forgeSchema := map[string]configs.ForgeTmpl{}
		err = yaml.Unmarshal(asset, &forgeSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal the forge schema: %s", err)
		}

		forgeTmpl, ok := forgeSchema[c.ForgeTmplKeyName]
		if !ok {
			return nil, fmt.Errorf("failed to get %s tmpl from forge schema: %s", c.ForgeTmplKeyName, err)

		}
		return &forgeTmpl, nil
	}
	return nil, nil
}

func (c *Cfg) InitGoModulesInOutPath() error {
	var err error

	if c.OutDir != "" {
		log.Printf("running go mod init %s", c.OutDir)
		err = execInPath(c.OutDir, "go", "mod", "init", c.GoModulesProjectName)
		if err != nil {
			return fmt.Errorf("failed to init go modules: %s", err)
		}

		log.Println("running go mod tidy")
		err = execInPath(c.OutDir, "go", "mod", "tidy")
		if err != nil {
			return fmt.Errorf("failed to tidy go modules: %s", err)
		}
	}

	return nil
}

func (c *Cfg) InitGitOrigin(path string) error {
	var err error

	if c.GitOrigin != "" {
		log.Println("git init")
		err = execInPath(path, "git", "init")
		if err != nil {
			return fmt.Errorf("failed to init git repository: %s", err)
		}

		log.Println("git add .")
		err = execInPath(path, "git", "add", ".")
		if err != nil {
			return fmt.Errorf("failed to add all chahnges to git repository: %s", err)
		}

		log.Printf("git add origin: %s", c.GitOrigin)
		err = execInPath(path, "git", "remote", "add", "origin", c.GitOrigin)
		if err != nil {
			return fmt.Errorf("failed to add remote origin %s: %s", c.GitOrigin, err)
		}
	}

	return nil
}

func execInPath(projectPath, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = projectPath
	return cmd.Run()
}

type genTypeSurvey struct {
	ProjectGenType string
	WithGitOrigin  string
}

func askProjectGenTypeQuestion() (*genTypeSurvey, error) {
	var (
		answers = new(genTypeSurvey)
		err     error
	)

	var qs = []*survey.Question{
		{
			Name: "projectGenType",
			Prompt: &survey.Select{
				Message: "Choose project type:",
				Options: []string{ProjectGenGoPathType, ProjectGenGoModulesType},
				Default: ProjectGenGoPathType,
			},
		},
		{
			Name: "withGitOrigin",
			Prompt: &survey.Select{
				Message: "Initialize git repository",
				Options: []string{"yes", "no"},
				Default: "no",
			},
		},
	}

	err = survey.Ask(qs, answers)
	if err != nil {
		return nil, fmt.Errorf("failed to get the survey answer: %s", err)
	}
	return answers, nil
}

type genGoModulesSurvey struct {
	OutDirPath  string
	ProjectName string
}

func askGenGoModulesSurvey() (*genGoModulesSurvey, error) {
	var (
		answers = new(genGoModulesSurvey)
		err     error
	)

	var qs = []*survey.Question{
		{
			Name: "outDirPath",
			Prompt: &survey.Input{
				Message: "Enter the output dir to generate project with Go Modules:",
				Default: "./scaffold",
			},
			//Validate: survey.Required,
		},
		{
			Name: "projectName",
			Prompt: &survey.Input{
				Message: "Enter Project name of project with Go Modules (ex. forge)",
				Default: "forge",
			},
			Validate: survey.Required,
		},
	}

	err = survey.Ask(qs, answers)
	if err != nil {
		return nil, fmt.Errorf("failed to get the survey answer: %s", err)
	}
	return answers, nil
}
