package configs

import (
	"fmt"
	"os"
)

type ModelConfig struct {
	BaseConfig
	TPath string
}

// Validate is an implementation of Validatable interface from ozzo-validation.
func (config *ModelConfig) Validate() error {
	if err := config.BaseConfig.Validate(); err != nil {
		return err
	}
	if config.TPath == "" {
		return fmt.Errorf("tmpl: must be specified")
	}

	_, err := os.Stat(config.TPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("tmpl: file is not exist")
	}
	if err != nil {
		return fmt.Errorf("tmpl: %s", err.Error())
	}
	return nil
}
