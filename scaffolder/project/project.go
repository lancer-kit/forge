package project

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/lancer-kit/forge/configs"
)

type Project struct {
	Cfg *configs.ScaffolderCfg
}

// Scaffold scaffolds the bindata file tmpl
func (p *Project) Scaffold() error {
	for _, fileName := range AssetNames() {
		file := filepath.Join(p.Cfg.OutPath, fileName)
		genPath := strings.TrimSuffix(file, filepath.Ext(file))

		schema := map[string]interface{}{
			"project_name": p.Cfg.ProjectName,
		}
		err := RestoreTemplate(genPath, fileName, schema)
		if err != nil {
			return fmt.Errorf("failed to restore tmpl: %s", err)
		}
	}
	return nil
}

func ExecuteTemplate(name string, data interface{}) (*bytes.Buffer, error) {
	asset, err := Asset(name)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("").Parse(string(asset))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func RestoreTemplate(path, name string, data interface{}) error {
	buf, err := ExecuteTemplate(name, data)
	if err != nil {
		return err
	}

	info, err := AssetInfo(name)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf.Bytes(), info.Mode())
	if err != nil {
		return err
	}

	return nil
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
	if p.Cfg.OutPath == "" {
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

		p.Cfg.OutPath = projectGoPath
	}

	err := p.Cfg.Validate()
	if err != nil {
		log.Println("outPath or projectName param is not set")
		os.Exit(1)

	}
	return p
}
