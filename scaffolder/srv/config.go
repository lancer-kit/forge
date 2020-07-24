package srv

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
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
		var genTypeSurvey = new(genTypeSurvey)

		genTypeSurvey, err = askProjectGenTypeQuestion()
		if err != nil {
			return nil, fmt.Errorf("failed to get the survey answer: %s", err)
		}
		if genTypeSurvey.WithGitOrigin == "yes" {
			surveyCfg.GitOrigin, err = askPromptInput("Enter Git origin address:", "")
		}

		switch genTypeSurvey.ProjectGenType {
		case ProjectGenGoPathType:
			surveyCfg.GoPathDomainName, err = askPromptInput("Enter Project domain name` "+
				"(ex. gitlab.com/team/project) to be created in GOPATH", "forge")
			if err != nil {
				return nil, fmt.Errorf("failed to get the survey answer: %s", err)

			}
			return surveyCfg, nil

		case ProjectGenGoModulesType:
			var genModulesSurvey = new(genGoModulesSurvey)
			genModulesSurvey, err = askGenGoModulesSurvey()
			if err != nil {
				return nil, fmt.Errorf("failed to get the survey answer: %s", err)
			}

			surveyCfg.OutDir = genModulesSurvey.OutDirPath
			surveyCfg.GoModulesProjectName = genModulesSurvey.ProjectName

			return surveyCfg, nil
		}

	case TmplTypeForge:

	}
	return surveyCfg, err
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
