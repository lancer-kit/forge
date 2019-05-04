package configs

import (
	"github.com/lancer-kit/forge/templates"
)

type EnumsConfig struct {
	BaseConfig
	TransformRule templates.TransformRule
	AddTypePrefix bool
}

func (config *EnumsConfig) Validate() error {
	if err := config.BaseConfig.Validate(); err != nil {
		return err
	}
	if err := config.TransformRule.Validate(); err != nil {
		return err
	}

	return nil
}
