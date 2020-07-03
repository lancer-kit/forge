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
	FlagDomain     = "d"
	FlagName       = "n"
	FlagOutputPath = "o"
	FlagGoModules  = "gomods"

	FlagSchemaPath = "schema"
	// Optional flags that are used to scaffold custom project with some
	// defined workers api/db/
	FlagAPIService          = "api"
	FlagDBService           = "db"
	FlagSimpleWorkerService = "base_uwe"
)

func NewProjectCmd() cli.Command {
	return cli.Command{
		Name:   "new",
		Usage:  "generate new project structure from template",
		Action: scaffoldAction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  FlagGoModules,
				Usage: "Initializes the go modules with module name in scaffold project",
			},
			&cli.StringFlag{
				Name:  FlagOutputPath,
				Usage: "Specifies output dir to scaffold the project",
				Value: "./out",
			},
			&cli.StringFlag{
				Name:  FlagDomain,
				Usage: "Specifies project scaffold domain",
			},
			&cli.StringFlag{
				Name:  FlagName,
				Usage: "Specifies project scaffold name",
				Value: "scaffold/project",
			},
			&cli.BoolFlag{
				Name:  FlagAPIService,
				Usage: "Specifies generation of optional API service logic",
			},
			&cli.BoolFlag{
				Name:  FlagDBService,
				Usage: "Specifies generation of optional DB service logic",
			},
			&cli.BoolFlag{
				Name:  FlagSimpleWorkerService,
				Usage: "Specifies generation of optional simple uwe worker logic",
			},
			&cli.StringFlag{
				Name:   FlagSchemaPath,
				Usage:  "Specifies the tmpl schema path",
				Hidden: true,
				Value:  "./scaffolder/templates/schema.yml",
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

	tmplModules := configs.ScaffoldTmplModules{
		configs.ScaffoldProjectNameKey: projectName,
		configs.ModuleKeyAPI:           c.Bool(FlagAPIService),
		configs.ModuleKeyDB:            c.Bool(FlagDBService),
		configs.ModuleKeySimpleWorker:  c.Bool(FlagSimpleWorkerService),
	}

	cfg := configs.ScaffolderCfg{
		OutPath:     c.String(FlagOutputPath),
		Schema:      project.ReadSchema(c.String(FlagSchemaPath)),
		ProjectName: projectName,
		TmplModules: tmplModules,
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
