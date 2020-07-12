package configs

import (
	"github.com/lancer-kit/forge/templates"
)

type EnumsConfig struct {
	BaseConfig
	TransformRule templates.TransformRule
	AddTypePrefix bool
}

// Validate is an implementation of Validatable interface from ozzo-validation.
func (config *EnumsConfig) Validate() error {
	if err := config.BaseConfig.Validate(); err != nil {
		return err
	}
	if err := config.TransformRule.Validate(); err != nil {
		return err
	}

	return nil
}
