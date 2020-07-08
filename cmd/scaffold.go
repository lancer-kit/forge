package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/urfave/cli"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/scaffolder/project"
)

const (
	FlagProjectDomain         = "domain"
	FlagProjectName           = "name"
	FlagProjectPathWithGoMods = "gomods"
	FlagGitOrigin             = "gitorigin"
)

func NewProjectCmd() cli.Command {
	return cli.Command{
		Name:   "new",
		Usage:  "generate new project structure from template",
		Action: scaffoldAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  FlagProjectPathWithGoMods + ", m",
				Usage: "`dir path` to init project with Go Modules from ",
			},
			&cli.StringFlag{
				Name:  FlagProjectDomain + ", d",
				Usage: "`project domain name`, (ex. github.com, gitlab.com)",
			},
			&cli.StringFlag{
				Name:  FlagProjectName + ", n",
				Usage: "`project name`",
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
	log.Println(cfg)

	scaffoldProject := project.NewProject(cfg)
	err = scaffoldProject.Scaffold()
	if err != nil {
		return fmt.Errorf("failed to scaffold project: %s", err)
	}

	projectPath := c.String(FlagProjectPathWithGoMods)
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
	return nil
}

func execInScaffoldPath(projectPath, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = projectPath
	return cmd.Run()
}

type RepositoryDomain string

const (
	DomainGitHub RepositoryDomain = "github.com"
	DomainGitLab RepositoryDomain = "gitlab.com"
)

type ScaffoldCliValues struct {
	dirPathWithGoMods string
	projectDomain     RepositoryDomain
	name              string
	gitOrigin         string
}

func (c ScaffoldCliValues) Validate() error {
	return validation.Errors{
		FlagProjectName:   validation.Validate(&c.name, validation.Required),
		FlagProjectDomain: validation.Validate(&c.projectDomain, validation.In(DomainGitHub, DomainGitLab)),
		FlagGitOrigin:     validation.Validate(&c.gitOrigin),
		FlagProjectPathWithGoMods: validation.Validate(c.dirPathWithGoMods,
			validation.Match(regexp.MustCompile(`^[^-].*`))),
	}.Filter()
}

func scaffoldConfig(c *cli.Context) (*configs.ScaffolderCfg, error) {
	flagsValues := &ScaffoldCliValues{
		dirPathWithGoMods: c.String(FlagProjectPathWithGoMods),
		projectDomain:     RepositoryDomain(c.String(FlagProjectDomain)),
		name:              c.String(FlagProjectName),
		gitOrigin:         c.String(FlagGitOrigin),
	}
	err := flagsValues.Validate()
	if err != nil {
		return nil, fmt.Errorf("wrong cli flag: %s", err)
	}

	cfg := new(configs.ScaffolderCfg)

	if flagsValues.projectDomain == "" {
		cfg.ProjectName = flagsValues.name
	} else {
		cfg.ProjectName = fmt.Sprintf("%s/%s", flagsValues.projectDomain, flagsValues.name)
	}

	if c.String(FlagProjectPathWithGoMods) != "" {
		cfg.ProjectPath = c.String(FlagProjectPathWithGoMods)
	}
	return cfg, nil
}
