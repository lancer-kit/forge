package templates

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sheb-gregor/goplater/parser"
)

type ModelSpec struct {
	Package    string
	TypeName   string
	TypeString string
	Fields     []Field
}

type Field struct {
	Name  string
	FType string
	Tags  map[string]string
}

func (spec *ModelSpec) Exec(tmpl *template.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, spec)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}
	return buf.String(), nil
}

func OpenTemplate(templatePath string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, errors.Wrap(err, "invalid template")
	}
	return tmpl, nil
}

func FigureOut(spec *parser.StructureSpec) *ModelSpec {
	result := strings.ToLower(spec.Name[:1])
	result += spec.Name[1:]

	s := ModelSpec{
		TypeName:   spec.Name,
		TypeString: result,
	}
	for _, fieldName := range spec.Fields {
		s.Fields = append(s.Fields, Field{
			Name:  fieldName,
			FType: spec.FTypes[fieldName],
			Tags:  parseTag(spec.Tags[fieldName]),
		})
	}
	return &s
}

func parseTag(rawTag string) map[string]string {
	rawTag = strings.Trim(rawTag, "`")
	tags := map[string]string{}
	for _, fullTag := range strings.Split(rawTag, " ") {
		tag := strings.Split(fullTag, ":")
		tags[tag[0]] = strings.Trim(strings.Split(tag[1], ",")[0], `"`)
	}
	return tags
}
