package main

import (
	"strings"

	"github.com/fatih/camelcase"
)

var Trasformers = map[string]func(string) string{
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
