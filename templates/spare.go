package templates

var Spare = map[string]CodeTemplate{
	"NameToValue": {Name: "NameToValue", Raw: nameToValueSpareRaw},
	"ValueToName": {Name: "ValueToName", Raw: valueToNameSpareRaw},
}

func init() {
	for i, t := range Spare {
		t.parse()
		Spare[i] = t
	}
}

var (
	nameToValueSpareRaw = `
var def{{.TypeName}}NameToValue = map[string]{{.TypeName}} {
        {{range .Values}}{{.Name}}.String(): {{.Name}},
        {{end}}
    }

`

	valueToNameSpareRaw = `
var def{{.TypeName}}ValueToName = map[{{.TypeName}}]string {
        {{range .Values}}{{.Name}}: {{.Name}}.String(),
        {{end}}
    }
`
)
