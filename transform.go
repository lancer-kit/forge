package main

import (
	"strings"

	"github.com/fatih/camelcase"
	"github.com/sheb-gregor/goplater/templates"
)

type TransformRule string

var (
	TransformRuleSnake TransformRule = "snake"
	TransformRuleKebab TransformRule = "kebab"
	TransformRuleSpace TransformRule = "space"
	TransformRuleNone  TransformRule = "none"
)

func (rule TransformRule) Transform(src string) string {
	switch rule {
	case TransformRuleSnake:
		return transformString(src, "_")
	case TransformRuleKebab:
		return transformString(src, "-")
	case TransformRuleSpace:
		return transformString(src, " ")
	case TransformRuleNone:
		return src
	}

	return src
}

func transformString(src, delimiter string) string {
	entries := camelcase.Split(src)
	if len(entries) <= 1 {
		return strings.ToLower(src)
	}

	result := strings.ToLower(entries[0])
	for i := 1; i < len(entries); i++ {
		result += delimiter + strings.ToLower(entries[i])
	}
	return result
}

func (rule TransformRule) TransformValues(typeName string, values []string, keepTPrefix bool) []templates.TypeValue {
	var str string
	res := make([]templates.TypeValue, len(values))

	for i := range values {
		str = values[i]
		if !keepTPrefix {
			str = strings.Replace(str, typeName, "", 1)
		}

		res[i] = templates.TypeValue{
			Name: values[i],
			Str:  rule.Transform(str),
		}
	}
	return res
}
