package templates

import (
	"bytes"
	"go/format"
	"log"
	"text/template"
)

type CodeTemplate struct {
	Name   string
	Raw    string
	Parsed *template.Template
}

func (r *CodeTemplate) parse() {
	r.Parsed = template.Must(template.New(r.Name).Parse(r.Raw))
}

type Analysis struct {
	Command     string
	PackageName string
	//TypesAndValues map[string][]TypeValue
	Types map[string]TypeSpec
}

type TypeSpec struct {
	TypeName    string
	Values      []TypeValue
	ExcludeList map[string]bool
}

type TypeValue struct {
	Name string
	Str  string
}

func (analysis *Analysis) GenerateByTemplate(merge bool) map[string][]byte {
	var results = make(map[string][]byte)

	var buf bytes.Buffer

	FileBase.Parsed.Execute(&buf, analysis)

	for typeName, spec := range analysis.Types {
		for _, t := range EnumBase {
			_, excludable := spec.ExcludeList[t.Name]
			//_, haveSpare := Spare[t.Name]
			if excludable {
				//if !haveSpare {
				continue
				//}
			}
			if err := t.Parsed.Execute(&buf, &spec); err != nil {
				log.Fatalf("generating code: %v", err)
			}
		}

		if !merge {
			results[typeName] = buf.Bytes()
			buf = bytes.Buffer{}
			FileBase.Parsed.Execute(&buf, analysis)
		}
	}
	if merge {
		results["all"] = buf.Bytes()
	}

	var err error
	for key, val := range results {
		results[key], err = format.Source(val)
		if err != nil {
			// Should never happen, but can arise when developing this code.
			// The user can compile the output to see the error.
			log.Printf("warning: internal error: invalid Go generated: %s", err)
			log.Printf("warning: compile the package to analyze the error")
			results[key] = val
		}
	}
	return results
}
