package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/urfave/cli"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/scaffolder/project"
)

const (
	FlagGoModsProjectName   = "gomods"
	FlagGoModsProjectPath   = "outdir"
	FlagProjectOriginGoPath = "gopath"
	FlagGitOrigin           = "gitorigin"

	CliMsgSuccess = "New project was successfully generated."
)

func NewProjectCmd() cli.Command {
	return cli.Command{
		Name:   "new",
		Usage:  "generate new project structure from template",
		Action: scaffoldAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  FlagGoModsProjectPath + ", o",
				Usage: "`dir path` to init project with Go Modules",
			},
			cli.StringFlag{
				Name:  FlagGoModsProjectName + ", m",
				Usage: "`project name` of project цшер Go Modules",
			},
			&cli.StringFlag{
				Name:  FlagProjectOriginGoPath + ", g",
				Usage: "`project domain name`",
			},
			&cli.StringFlag{
				Name:  FlagGitOrigin + ", r",
				Usage: "`git origin` to init git repository add all changes to remote origin",
			},
		},
	}
}

func scaffoldAction(c *cli.Context) error {
	cfg, err := scaffoldConfig(c)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to parse scaffold cmd flags: %s", err)
	}

	scaffoldProject := project.NewProject(cfg)
	err = scaffoldProject.Scaffold()
	if err != nil {
		return fmt.Errorf("failed to scaffold project: %s", err)
	}

	projectPath := c.String(FlagGoModsProjectPath)
	if projectPath != "" {
		log.Printf("running go mod init %s", cfg.ProjectName)
		err = execInScaffoldPath(projectPath, "go", "mod", "init", cfg.ProjectName)
		if err != nil {
			return fmt.Errorf("failed to init go modules: %s", err)
		}

		log.Println("running go mod tidy")
		err = execInScaffoldPath(projectPath, "go", "mod", "tidy")
		if err != nil {
			return fmt.Errorf("failed to tidy go modules: %s", err)
		}
	}
	if c.String(FlagGitOrigin) != "" {
		log.Println("git init")
		err = execInScaffoldPath(scaffoldProject.Cfg.ProjectPath, "git", "init")
		if err != nil {
			return fmt.Errorf("failed to init git repository: %s", err)
		}

		log.Println("git add .")
		err = execInScaffoldPath(scaffoldProject.Cfg.ProjectPath, "git", "add", ".")
		if err != nil {
			return fmt.Errorf("failed to add all chahnges to git repository: %s", err)
		}

		log.Printf("git add origin: %s", c.String(FlagGitOrigin))
		err = execInScaffoldPath(scaffoldProject.Cfg.ProjectPath, "git", "remote", "add", "origin", c.String(FlagGitOrigin))
		if err != nil {
			return fmt.Errorf("failed to add remote origin %s: %s", c.String(FlagGitOrigin), err)
		}
	}

	path, err := filepath.Abs(scaffoldProject.Cfg.ProjectPath)
	if err != nil {
		return err
	}
	log.Printf("%s Generated project [%s] in path: [%s]", CliMsgSuccess, scaffoldProject.Cfg.ProjectName, path)
	return nil
}

func execInScaffoldPath(projectPath, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = projectPath
	return cmd.Run()
}

type ScaffoldCliValues struct {
	projectPathWithGoMods string
	projectNameWithGoMods string
	projectGoPathOrigin   string
	gitOrigin             string
}

func (c ScaffoldCliValues) Validate() error {
	if c.projectNameWithGoMods == "" && c.projectGoPathOrigin == "" {
		return fmt.Errorf("specify the way of project generation gomod(--%s --%s flags) (gopath --%s flag)",
			FlagGoModsProjectPath, FlagGoModsProjectName, FlagProjectOriginGoPath)
	}
	return validation.Errors{
		FlagGoModsProjectName: validation.Validate(&c.projectNameWithGoMods, validation.When(
			c.projectPathWithGoMods != "", validation.Required,
			validation.Match(regexp.MustCompile(`^[^-].*`))),
		),
		FlagGoModsProjectPath: validation.Validate(&c.projectPathWithGoMods, validation.When(
			c.projectNameWithGoMods != "", validation.Required),
		),
		FlagGitOrigin:           validation.Validate(&c.gitOrigin),
		FlagProjectOriginGoPath: validation.Validate(&c.gitOrigin),
	}.Filter()
}

func scaffoldConfig(c *cli.Context) (*configs.ScaffolderCfg, error) {
	flagsValues := &ScaffoldCliValues{
		projectPathWithGoMods: c.String(FlagGoModsProjectPath),
		projectNameWithGoMods: c.String(FlagGoModsProjectName),
		projectGoPathOrigin:   c.String(FlagProjectOriginGoPath),
		gitOrigin:             c.String(FlagGitOrigin),
	}
	err := flagsValues.Validate()
	if err != nil {
		return nil, fmt.Errorf("cli error: %s", err)
	}

	cfg := new(configs.ScaffolderCfg)

	if flagsValues.projectNameWithGoMods != "" {
		cfg.ProjectName = flagsValues.projectNameWithGoMods
		cfg.ProjectPath = flagsValues.projectPathWithGoMods
	} else {
		cfg.ProjectName = flagsValues.projectGoPathOrigin
	}
	return cfg, nil
}
