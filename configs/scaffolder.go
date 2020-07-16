package configs

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const ForgeSchemaAssetName = "schema.yml"

type ScaffolderCfg struct {
	ProjectPath string
	ProjectName string

	ForgeTmplKeyName string
	ForgeTmpl        *ForgeTmpl
}

// Validate is an implementation of Validatable interface from ozzo-validation.
func (cfg ScaffolderCfg) Validate() error {
	return validation.ValidateStruct(&cfg,
		validation.Field(&cfg.ProjectPath, validation.Required),
		validation.Field(&cfg.ProjectName, validation.Required),
	)
}

type ForgeSchema map[string]ForgeTmpl

type ForgeTmpl struct {
	AssetPrefix string                 `yml:"assetprefix"`
	Fields      map[string]interface{} `yml:"fields"`
}

func (t ForgeTmpl) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.AssetPrefix, validation.Required),
		validation.Field(&t.Fields, validation.Required),
	)
}
