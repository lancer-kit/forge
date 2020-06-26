package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/urfave/cli"

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

func ScaffoldCmd() cli.Command {
	cmd := cli.Command{
		Name:  "gen",
		Usage: "gen scaffold project structure",
		Action: func(c *cli.Context) error {
			schema := project.ReadSchema(c.String(FlagSchemaPath))
			var projectName string
			if c.String(FlagDomain) == "" {
				projectName = c.String(FlagName)
			} else {
				projectName = fmt.Sprintf("%s/%s", c.String(FlagDomain), c.String(FlagName))
			}
			log.Printf("scaffolding project %s", projectName)

			scaffoldSchema := project.ScaffoldTmplModules{
				project.ScaffoldProjectNameKey: projectName,
				project.ModuleKeyAPI:           c.Bool(FlagAPIService),
				project.ModuleKeyDB:            c.Bool(FlagDBService),
				project.ModuleKeySimpleWorker:  c.Bool(FlagSimpleWorkerService),
			}
			p := project.NewProject(c.String(FlagOutputPath), schema, scaffoldSchema)

			err := p.Scaffold()
			if err != nil {
				return fmt.Errorf("failed to scaffold project: %s", err)
			}

			if c.Bool(FlagGoModules) {
				log.Printf("running go mod init %s", projectName)
				err = execInScaffoldPath(c, "go", "mod", "init", projectName)
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
		},
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
				Value:  "../scaffolder/templates/schema.yml",
			},
		},
	}
	return cmd
}

func execInScaffoldPath(c *cli.Context, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = c.String(FlagOutputPath)
	return cmd.Run()
}
