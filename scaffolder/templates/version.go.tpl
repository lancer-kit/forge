package main

import (
	"{{.project_name}}/info"
)

var (
	Version = "1.0.0-rc"
	Build   string
	Tag     string
)

func init() {
	info.App.Version = Version
	info.App.Build = Build
	info.App.Tag = Tag
}
