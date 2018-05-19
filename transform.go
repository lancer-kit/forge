package main

import (
	"errors"
	"strings"

	"github.com/fatih/camelcase"
)

type TypeValue struct {
	Name string
	Str  string
}

var transformers = map[string]func(string) string{
	"snake": func(src string) string {
		return transformString(src, "_")
	},
	"kebab": func(src string) string {
		return transformString(src, "-")
	},
	"space": func(src string) string {
		return transformString(src, " ")
	},
	"none": func(src string) string {
		return src
	},
}

func transformString(src, delim string) string {
	entries := camelcase.Split(src)
	if len(entries) <= 1 {
		return strings.ToLower(src)
	}

	result := strings.ToLower(entries[0])
	for i := 1; i < len(entries); i++ {
		result += delim + strings.ToLower(entries[i])
	}
	return result
}

func transformValues(typeName string, values []string) ([]TypeValue, error) {
	if transformMethod == nil {
		return nil, errors.New("transform method is not defined")
	}

	transform, ok := transformers[*transformMethod]
	if !ok {
		return nil, errors.New("invalid transform method")
	}

	var str string
	res := make([]TypeValue, len(values))

	for i := range values {
		str = values[i]
		if !*addTypePrefix {
			str = strings.Replace(str, typeName, "", 1)
		}

		res[i] = TypeValue{
			Name: values[i],
			Str:  transform(str),
		}
	}
	return res, nil
}
