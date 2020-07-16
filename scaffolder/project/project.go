package project

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/scaffolder/project/provider"
)

type Project struct {
	Provider provider.BindataProvider
	Cfg      *configs.ScaffolderCfg
}

// Scaffold scaffolds the bindata file tmpl
func (p *Project) Scaffold() error {
	return p.Provider.Scaffold(p.Cfg.ProjectPath, p.Cfg.ProjectName)
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}

	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}

	return path[0:len(prefix)] == prefix
}

func NewProject(cfg *configs.ScaffolderCfg) *Project {
	p := new(Project)
	p.Cfg = cfg

	// Check if the outPath is empty(in case if Go Modules are disabled and project will be generated
	// into $GOPATH) or not (in case if Go Modules are enabled and project will be generated by the --out path flag
	if p.Cfg.ProjectPath == "" {
		// find the $GOPATH to generate the project directory
		envGoPath := os.Getenv("GOPATH")
		goPaths := filepath.SplitList(envGoPath)
		if len(goPaths) == 0 {
			log.Println("$GOPATH is not set")
			os.Exit(1)
		}
		srcPaths := make([]string, 0, len(goPaths))
		for _, goPath := range goPaths {
			srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
		}

		// in case if the user in GOPATH generate project
		var projectGoPath string
		wd, err := os.Getwd()
		if err != nil {
			return nil
		}
		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) {
				projectGoPath = filepath.Join(srcPath, p.Cfg.ProjectName)
				break
			}
		}

		// in case if user is not in GOPATH generate project
		if projectGoPath == "" {
			projectGoPath = filepath.Join(srcPaths[0], p.Cfg.ProjectName)
		}

		p.Cfg.ProjectPath = projectGoPath
	}

	err := p.Cfg.Validate()
	if err != nil {
		log.Println("outPath or projectName param is not set")
		os.Exit(1)
	}

	// set bindata provider
	if p.Cfg.ForgeTmplKeyName != "" {
		p.Provider = provider.NewBindataProvider(provider.Forge, p.Cfg.ForgeTmpl.AssetPrefix)
	} else {
		p.Provider = provider.NewBindataProvider(provider.Default, "")
	}

	log.Printf("path %s", p.Cfg.ProjectPath)
	return p
}
