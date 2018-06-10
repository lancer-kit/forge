package parser

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"strings"
)

// ValuesOfType is inspect files for constant value, default variable and methods for type,
// return a list of the constant values, and map of templates which must be ignored,
// because they have already been declared.
func (pkg *Package) ValuesOfType(typeName string) ([]string, map[string]bool, error) {
	var values, inspectErrs []string
	tmplsToExclude := map[string]bool{}

	for _, file := range pkg.files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch decl := node.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.CONST:
					vs, err := pkg.constOfTypeIn(typeName, decl)
					values = append(values, vs...)
					if err != nil {
						inspectErrs = append(inspectErrs, err.Error())
					}

				case token.VAR:
					vs := pkg.varOfTypeIn(typeName, decl)
					for k, v := range vs {
						tmplsToExclude[k] = v
					}
				default:
					return true
				}

			case *ast.FuncDecl:
				vs := pkg.methodsOfTypeIn(typeName, decl)
				for k, v := range vs {
					tmplsToExclude[k] = v
				}
			default:
				return true
			}

			return false
		})
	}

	if len(inspectErrs) > 0 {
		return nil, nil, fmt.Errorf("inspecting code:\n\t%v", strings.Join(inspectErrs, "\n\t"))
	}
	if len(values) == 0 {
		return nil, nil, fmt.Errorf("no values defined for type %s", typeName)
	}

	return values, tmplsToExclude, nil
}

// constOfTypeIn checks if a constant values is declared
// for the type and add it to the result list.
func (pkg *Package) constOfTypeIn(typeName string, decl *ast.GenDecl) ([]string, error) {
	var values []string

	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value, the constant is untyped.
			// Skip this vspec and reset the remembered type.
			typ = ""
			continue
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != typeName {
			// This is not the type we're looking for.
			continue
		}

		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := pkg.defs[name]
			if !ok {
				return nil, fmt.Errorf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				return nil, fmt.Errorf("can't handle non-integer constant type %s", typ)
			}
			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
			}
			values = append(values, name.Name)
		}
	}
	return values, nil
}

// varOfTypeIn  checks if a default variable is declared for the type,
// if declared - add it to the ignore list, and the template for this
// variable will NOT be added to the output file.
func (pkg *Package) varOfTypeIn(typeName string, decl *ast.GenDecl) map[string]bool {
	tmpls := map[string]bool{}

	for _, spec := range decl.Specs {
		vspec, ok := spec.(*ast.ValueSpec)
		if !ok {
			return tmpls
		}

		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}

			for _, v := range typeVariables {
				if strings.Contains(name.Name, v) && strings.Contains(name.Name, typeName) {
					tmpls[v] = true
				}

			}
		}
	}
	return tmpls
}
