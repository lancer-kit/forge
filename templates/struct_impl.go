package templates

type StructSpec struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
	Tags map[string]string
}
