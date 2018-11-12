package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/pkg/errors"
	"gitlab.inn4science.com/gophers/goplater/parser"
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

func FigureOut(spec *parser.StructureSpec) (_ *ModelSpec, err error) {
	result := strings.ToLower(spec.Name[:1])
	result += spec.Name[1:]

	s := ModelSpec{
		TypeName:   spec.Name,
		TypeString: result,
	}

	var tags map[string]string
	for _, fieldName := range spec.Fields {
		tags, err = parseTag(spec.Tags[fieldName])
		if err != nil {
			return nil, err
		}
		s.Fields = append(s.Fields, Field{
			Name:  fieldName,
			FType: spec.FTypes[fieldName],
			Tags:  tags,
		})
	}

	return &s, nil
}

func parseTag(rawTag string) (map[string]string, error) {
	tags := map[string]string{}

	// remove leading and trailing backtick
	rawTag = strings.Trim(rawTag, "`")

	for _, fullTag := range strings.Split(rawTag, " ") {
		// minimal valid tag is `k:"v"`
		// should contain:
		//    k - key name
		//    : - separator
		//    " - double quotes
		//    v - value
		// --> totally 5 chars
		if len(fullTag) < 5 {
			continue
		}

		// remove leading and trailing whitespaces, if exist
		fullTag = strings.Trim(fullTag, " ")

		tag := strings.Split(fullTag, ":")
		if len(tag) != 2 {
			return nil, fmt.Errorf("tag is invalid: %s", fullTag)
		}

		//tags can have not only values, but also optional parameters,
		// such as ʻomitempty' for JSON, separated by comma;
		//therefore, we divide the value by comma and take only the first value
		tagValue := strings.Split(tag[1], ",")[0]

		// remove leading and trailing double quotes, if exist
		tags[tag[0]] = strings.Trim(tagValue, `"`)
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("tag is invalid: %s", rawTag)
	}
	return tags, nil
}
