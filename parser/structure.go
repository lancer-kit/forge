package parser

import (
	"go/ast"
	"go/token"
)

type StructureSpec struct {
	Name   string
	Fields []string
	FTypes map[string]string
	Tags   map[string]string
}

// StructureDef is inspect files for the declaration of target structure type ...TODO
func (pkg *Package) FindStructureSpec(typeName string) (result *StructureSpec, err error) {
	for _, file := range pkg.files {
		if result != nil {
			return
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch decl := node.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.TYPE:
					for _, spec := range decl.Specs {
						res, ok := parseStructSpec(typeName, spec)
						if ok {
							result = res
							return false
						}
					}
				default:
					return true
				}

			default:
				return true
			}

			return true
		})
	}

	return
}

func parseStructSpec(typeName string, spec ast.Spec) (*StructureSpec, bool) {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, false
	}

	if typeSpec.Name.Name != typeName {
		return nil, false
	}

	structSpec, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		// todo: set error
		return nil, false
	}

	res := StructureSpec{
		Name:   typeName,
		FTypes: map[string]string{},
		Tags:   map[string]string{},
	}

	for _, f := range structSpec.Fields.List {
		ident, ok := f.Type.(*ast.Ident)
		if !ok {
			continue
		}

		// todo: add support of multi-name field declarations
		name := f.Names[0].Name

		res.Fields = append(res.Fields, name)
		res.Tags[name] = f.Tag.Value
		res.FTypes[name] = ident.Name

	}
	return &res, true
}
