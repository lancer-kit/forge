package templates

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"gitlab.inn4science.com/gophers/forge/parser"
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

	var tagsKV map[string]string
	for _, fieldName := range spec.Fields {
		tagsKV, err = parseTag(spec.Tags[fieldName])
		if err != nil {
			return nil, err
		}
		s.Fields = append(s.Fields, Field{
			Name:  fieldName,
			FType: spec.FTypes[fieldName],
			Tags:  tagsKV,
		})
	}

	return &s, nil
}

func parseTag(rawTag string) (map[string]string, error) {
	tags, err := parseRawTags(rawTag)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, errors.New("No tags")
	}

	tags = sanitizeTags(tags)
	return tags, nil
}

func parseRawTags(tag string) (map[string]string, error) {
	tags := map[string]string{}

	// remove leading and trailing backtick
	tag = strings.Trim(tag, "`")
	const kvSeparator = ':'
	const quote = '"'
	const whitespace = ' '

	var keyFound bool
	var keyStart, valueStart int
	var key string

	//todo: add comments
	for i := range tag {
		if keyFound && valueStart == i {
			continue
		}

		s := tag[i]
		if i < 1 && s == kvSeparator {
			// separator can not be in first position
			return nil, errors.New("invalid tag")
		}

		if !keyFound && s == kvSeparator {
			if tag[i+1] != quote {
				return nil, errors.New("invalid tag")
			}
			key = tag[keyStart:i]
			invalid, _ := regexp.MatchString("([^a-zA-Z0-9_]+)", key)
			if invalid {
				return nil, errors.New("invalid key")
			}
			valueStart = i + 1
			keyFound = true
			continue
		}

		if keyFound && tag[i] == quote {
			if i+1 < len(tag) && tag[i+1] != whitespace {
				return nil, errors.New("invalid tag")
			}
			// remove leading and trailing double quotes, if exist
			tags[key] = strings.Trim(tag[valueStart:i], `"`)
			keyStart = i + 2
			keyFound = false
			continue
		}

	}
	return tags, nil
}

func sanitizeTags(tags map[string]string) map[string]string {
	for key, value := range tags {
		// remove leading and trailing double quotes, if exist
		value = strings.Trim(value, `"`)
		// tags can have not only values, but also optional parameters,
		// such as 'omitempty' for JSON, separated by comma;
		// therefore, we divide the value by comma and take only the first value
		value = strings.Split(value, ",")[0]
		// special case for hidden fields
		if value == "-" {
			delete(tags, key)
			continue
		}
		tags[key] = value
	}
	return tags
}
