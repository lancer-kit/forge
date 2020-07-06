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
	FlagDomain     = "domain"
	FlagName       = "name"
	FlagOutputPath = "output"
	FlagGoModules  = "gomods"
)

func NewProjectCmd() cli.Command {
	return cli.Command{
		Name:   "new",
		Usage:  "generate new project structure from template",
		Action: scaffoldAction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  FlagGoModules + ", m",
				Usage: "Initializes the go modules with module name in scaffold project",
			},
			&cli.StringFlag{
				Name:  FlagOutputPath + ", o",
				Usage: "Specifies output dir to scaffold the project",
				Value: "./out",
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
		},
	}
}

func scaffoldAction(c *cli.Context) error {
	cfg := scaffoldConfig(c)

	err := project.NewProject(&cfg).Scaffold()
	if err != nil {
		return fmt.Errorf("failed to scaffold project: %s", err)
	}

	if c.Bool(FlagGoModules) {
		log.Printf("running go mod init %s", cfg.ProjectName)
		err = execInScaffoldPath(c, "go", "mod", "init", cfg.ProjectName)
		if err != nil {
			return fmt.Errorf("failed to init go modules: %s", err)
		}

		log.Println("running go mod tidy")
		err = execInScaffoldPath(c, "go", "mod", "tidy")
		if err != nil {
			return fmt.Errorf("failed to tidy go modules: %s", err)
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

	cfg := configs.ScaffolderCfg{
		OutPath:     c.String(FlagOutputPath),
		ProjectName: projectName,
	}
	err := cfg.Validate()
	if err != nil {
		log.Fatalf("no all necessary fields for scaffolding the project: %s", err)
	}
	return cfg
}

func execInScaffoldPath(c *cli.Context, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = c.String(FlagOutputPath)
	return cmd.Run()
}
