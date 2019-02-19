package configs

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

type BaseConfig struct {
	Types        []string
	MergeSpecs   bool
	OutputDir    string
	OutputName   string
	OutputSuffix string
	OutputPrefix string
}

func (config *BaseConfig) Validate() error {
	if len(config.Types) == 0 {
		return fmt.Errorf("type: should not be empty")
	}
	if config.OutputPrefix == "" && config.OutputSuffix == "" {
		return fmt.Errorf("sufix or prefix: should be passed")
	}

	return nil
}

func (config BaseConfig) GetPath(name, dir string) string {
	var splittedName string

	for i, r := range name {
		if i == 0 {
			splittedName += string(r)
			continue
		}
		if unicode.IsUpper(r) {
			splittedName += "_" + string(r)
			continue
		}
		splittedName += string(r)
	}

	output := strings.ToLower(config.OutputPrefix + splittedName + config.OutputSuffix + ".go")

	return filepath.Join(dir, output)
}
