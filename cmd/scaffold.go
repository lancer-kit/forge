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

	//var forgeTmplKeyName = flagsValues.withForgeTmpl
	//if forgeTmplKeyName != "" {
	//	cfg.ForgeTmplKeyName = forgeTmplKeyName
	//
	//	// Get Forge schema to config
	//	asset, err := forge.Asset(configs.ForgeSchemaAssetName)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to load the %s asset: %s", configs.ForgeSchemaAssetName, err)
	//	}
	//
	//	forgeSchema := map[string]configs.ForgeTmpl{}
	//	err = yaml.Unmarshal(asset, &forgeSchema)
	//	if err != nil {
	//		log.Println(err)
	//
	//		return nil, fmt.Errorf("failed to unmarshal the forge schema: %s", err)
	//	}
	//	forgeTmpl, ok := forgeSchema[forgeTmplKeyName]
	//	if !ok {
	//		return nil, fmt.Errorf("failed to get %s tmpl from forge schema: %s", forgeTmplKeyName, err)
	//
	//	}
	//	cfg.ForgeTmpl = &forgeTmpl
	//}
	return cfg, nil
}
