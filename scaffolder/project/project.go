package project

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lancer-kit/forge/configs"
)

type Project struct {
	cfg *configs.ScaffolderCfg
}

// Scaffold scaffolds the bindata file tmpl
func (p *Project) Scaffold() error {
	for _, fileName := range AssetNames() {
		file := filepath.Join(p.cfg.OutPath, fileName)
		genPath := strings.TrimSuffix(file, filepath.Ext(file))

		schema := map[string]interface{}{
			"project_name": p.cfg.ProjectName,
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

func NewProject(cfg *configs.ScaffolderCfg) *Project {
	return &Project{
		cfg: cfg,
	}
}
