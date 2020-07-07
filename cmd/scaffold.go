package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/urfave/cli"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/scaffolder/project"
)

const (
	FlagDomain               = "domain"
	FlagName                 = "name"
	FlagGoModulesProjectPath = "gomods"
	FlagGitInitRepo          = "repo"
)

func NewProjectCmd() cli.Command {
	return cli.Command{
		Name:   "new",
		Usage:  "generate new project structure from template",
		Action: scaffoldAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  FlagGoModulesProjectPath + ", m",
				Usage: "Initializes the go modules with module name in scaffold project",
			},
			&cli.StringFlag{
				Name:  FlagDomain + ", d",
				Usage: "Specifies project scaffold domain",
			},
			&cli.StringFlag{
				Name:  FlagName + ", n",
				Usage: "Specifies project scaffold name",
				Value: "scaffold/project",
			},
			&cli.StringFlag{
				Name:  FlagGitInitRepo + ", r",
				Usage: "Initialize git repository with origin",
			},
		},
	}
}

func scaffoldAction(c *cli.Context) error {
	cfg := scaffoldConfig(c)

	scaffoldProject := project.NewProject(&cfg)
	err := scaffoldProject.Scaffold()
	if err != nil {
		return fmt.Errorf("failed to scaffold project: %s", err)
	}

	projPath := c.String(FlagGoModulesProjectPath)
	if projPath != "" {
		log.Printf("running go mod init %s", cfg.ProjectName)
		err = execInScaffoldPath(projPath, "go", "mod", "init", cfg.ProjectName)
		if err != nil {
			return fmt.Errorf("failed to init go modules: %s", err)
		}

		log.Println("running go mod tidy")
		err = execInScaffoldPath(projPath, "go", "mod", "tidy")
		if err != nil {
			return fmt.Errorf("failed to tidy go modules: %s", err)
		}
	}

	if c.String(FlagGitInitRepo) != "" {
		log.Println("git init")
		err = execInScaffoldPath(scaffoldProject.Cfg.OutPath, "git", "init")
		if err != nil {
			return fmt.Errorf("failed to init git repository: %s", err)
		}

		log.Println("git add .")
		err = execInScaffoldPath(scaffoldProject.Cfg.OutPath, "git", "add", ".")
		if err != nil {
			return fmt.Errorf("failed to add all chahnges to git repository: %s", err)
		}

		log.Printf("git add origin: %s", c.String(FlagGitInitRepo))
		err = execInScaffoldPath(scaffoldProject.Cfg.OutPath, "git", "remote", "add", "origin", c.String(FlagGitInitRepo))
		if err != nil {
			return fmt.Errorf("failed to add remote origin %s: %s", c.String(FlagGitInitRepo), err)
		}
	}
	return nil
}

func scaffoldConfig(c *cli.Context) configs.ScaffolderCfg {
	var projectName string
	if c.String(FlagDomain) == "" {
		projectName = c.String(FlagName)
	} else {
		projectName = fmt.Sprintf("%s/%s", c.String(FlagDomain), c.String(FlagName))
	}

	var projectPath string
	if c.String(FlagGoModulesProjectPath) != "" {
		projectPath = c.String(FlagGoModulesProjectPath)
	}

	cfg := configs.ScaffolderCfg{
		OutPath:     projectPath,
		ProjectName: projectName,
	}
	return cfg
}

func execInScaffoldPath(projectPath, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = projectPath
	return cmd.Run()
}
