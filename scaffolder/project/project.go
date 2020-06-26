package project

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Project struct {
	outPath    string
	tmplData   ScaffoldTmplModules
	modulesCfg TemplatesCfg
}

// Scaffold scaffolds the bindata file tmpl
func (p *Project) Scaffold() error {
	// generate the base directory structure
	err := p.scaffoldTmplDir(p.modulesCfg.Target.Path)
	if err != nil {
		return fmt.Errorf("failed to scaffold base dir: %s", err)
	}

	for tmplKey := range p.tmplData {
		key, ok := tmplKey.(ScaffoldTmplKey)
		if !ok {
			continue
		}
		if p.tmplData[tmplKey] == true {
			log.Printf("generate the template %s %v", tmplKey, p.tmplData[tmplKey])
			module := p.modulesCfg.Modules[key]
			err := p.scaffoldTmplDir(module.Path)
			if err != nil {
				return fmt.Errorf("failed to scaffold base dir: %s", err)
			}
		}
	}
	return nil
}

func (p *Project) scaffoldTmplDir(dir string) error {
	assets := getAssetFromDir(dir)
	for _, fileName := range assets {
		relPath, err := filepath.Rel(dir, fileName)
		if err != nil {
			return fmt.Errorf("failed to get rel path: %s", err)
		}
		genRawPath := filepath.Join(p.outPath, relPath)
		genPath := strings.TrimSuffix(genRawPath, filepath.Ext(genRawPath))

		err = RestoreTemplate(genPath, fileName, p.tmplData)
		if err != nil {
			return fmt.Errorf("failed to restore tmpl: %s", err)
		}
	}
	return nil
}

func getAssetFromDir(dir string) []string {
	var paths = make([]string, 0)
	for _, tmplName := range AssetNames() {
		assetRootDirName := strings.Split(tmplName, "/")[0]
		if assetRootDirName == "" {
			return nil
		}
		if dir == filepath.Base(assetRootDirName) {
			paths = append(paths, tmplName)
		}
	}
	return paths
}

func executeTemplate(name string, data interface{}) (*bytes.Buffer, error) {
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
	buf, err := executeTemplate(name, data)
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

func NewProject(outPath string, cfg TemplatesCfg, data ScaffoldTmplModules) *Project {
	return &Project{
		outPath:    outPath,
		modulesCfg: cfg,
		tmplData:   data,
	}
}
