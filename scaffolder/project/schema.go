package project

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/lancer-kit/forge/configs"
)

func ReadSchema(path string) configs.TemplatesCfg {
	rawConfig, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read scaffold schema config file with path %s: %s", path, err)
	}

	config := new(configs.TemplatesCfg)
	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		log.Fatalf("unable to scaffold schema config file with raw config %s: %s", rawConfig, err)
	}

	err = config.Validate()
	if err != nil {
		log.Fatalf("invalid scaffold schema config: %s", err)
	}
	return *config
}
