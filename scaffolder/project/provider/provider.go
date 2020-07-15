package provider

import (
	"github.com/lancer-kit/forge/scaffolder/project/provider/defaultp"
	"github.com/lancer-kit/forge/scaffolder/project/provider/forge"
)

// Name represents Bindata provider name
type Name string

const (
	Default Name = "default"
	Forge   Name = "forge"
)

// BindataProvider represents provider for Bindata either it will
// be .forge or embedded provider from /templates project dir
type BindataProvider interface {
	Scaffold(outPath, projectName string) error
	GenerateTemplates(name string, path string, data interface{}) error
}

func NewBindataProvider(provider Name) BindataProvider {
	switch provider {
	case Forge:
		return &forge.ForgeProvider{}
	case Default:
		return &defaultp.DefaultProvider{}
	}
	return nil
}
