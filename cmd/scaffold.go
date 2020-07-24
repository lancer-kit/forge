package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/scaffolder/project"
	"github.com/lancer-kit/forge/scaffolder/srv"
)

const (
	CliMsgSuccess = "New project was successfully generated."
)

func NewProjectCmd() cli.Command {
	return cli.Command{
		Name:   "new",
		Usage:  "generate new project structure from template",
		Action: scaffoldAction,
	}
}

func scaffoldAction(c *cli.Context) error {
	cfg, err := srv.AskSurvey()
	if err != nil {
		return err
	}

	scaffoldCfg, err := scaffoldConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to parse scaffold cmd flags: %s", err)
	}

	scaffoldProject := project.NewProject(scaffoldCfg)
	err = scaffoldProject.Scaffold()
	if err != nil {
		return fmt.Errorf("failed to scaffold project: %s", err)
	}

	err = cfg.InitGoModulesInOutPath()
	if err != nil {
		return err
	}

	err = cfg.InitGitOrigin(scaffoldProject.Cfg.ProjectPath)
	if err != nil {
		return err
	}

	path, err := filepath.Abs(scaffoldProject.Cfg.ProjectPath)
	if err != nil {
		return err
	}
	log.Printf("%s Generated project [%s] in path: [%s]", CliMsgSuccess, scaffoldProject.Cfg.ProjectName, path)
	return nil
}

func scaffoldConfig(srvCfg *srv.Cfg) (*configs.ScaffolderCfg, error) {
	cfg := new(configs.ScaffolderCfg)

	if srvCfg.GoModulesProjectName != "" {
		cfg.ProjectName = srvCfg.GoModulesProjectName
		cfg.ProjectPath = srvCfg.OutDir
	} else {
		cfg.ProjectName = srvCfg.GoPathDomainName
	}

	tmpl, err := srvCfg.GetForgeTmpl()
	if err != nil {
		return nil, err
	}
	cfg.ForgeTmpl = tmpl

	return cfg, nil
}
