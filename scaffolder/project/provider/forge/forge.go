package forge

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	forge "github.com/lancer-kit/forge/.forge"
	"github.com/lancer-kit/forge/configs"
)

type ForgeProvider struct {
	AssetNamePrefix string
}

func (p *ForgeProvider) Scaffold(outPath, projectName string) error {
	for _, fileName := range forge.AssetNames() {
		if filepath.Ext(fileName) == filepath.Ext(configs.ForgeSchemaAssetName) {
			continue
		}
		fileNameNoPrefix, err := filepath.Rel(p.AssetNamePrefix, fileName)
		if err != nil {
			log.Println(err)
			return err
		}

		file := filepath.Join(outPath, fileNameNoPrefix)
		genPath := strings.TrimSuffix(file, filepath.Ext(file))

		schema := map[string]interface{}{
			"project_name": projectName,
		}

		err = p.GenerateTemplates(fileName, genPath, schema)
		if err != nil {
			return fmt.Errorf("failed to restore tmpl: %s", err)
		}
	}
	return nil
}

func (p *ForgeProvider) GenerateTemplates(name string, path string, data interface{}) error {
	asset, err := forge.Asset(name)
	if err != nil {
		return err
	}

	info, err := forge.AssetInfo(name)
	if err != nil {
		return err
	}

	err = generateTemplates(asset, info, path, data)
	if err != nil {
		return err
	}
	return nil
}

func generateTemplates(asset []byte, assetInfo os.FileInfo, path string, data interface{}) error {
	buf, err := executeTemplate(asset, data)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf.Bytes(), assetInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

func executeTemplate(asset []byte, data interface{}) (*bytes.Buffer, error) {
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
