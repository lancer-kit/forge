// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// Package parser parses Go code and keeps track of all the types defined
// and provides access to all the constants defined for an int type.
package parser

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

// typeVariables is slice of default variable for a type
// which will be generated from template.
var typeVariables = []string{
	"NameToValue",
	"ValueToName",
}

// typeMethods is a map of default methods for a type,
// which will be generated from template:
//    key - string name of the method,
//    value - show is receiver should be pointer or not.
// map [methodName]shouldBePointer
var typeMethods = map[string]bool{
	"String":        false,
	"Validate":      false,
	"MarshalJSON":   false,
	"UnmarshalJSON": true,
	"Value":         false,
	"Scan":          true,
}

// A Package contains all the information related to a parsed package.
type Package struct {
	Name  string
	files []*ast.File

	defs map[*ast.Ident]types.Object
}

// ParsePackage parses the package in the given directory and returns it.
func ParsePackage(directory string) (*Package, error) {
	relDir, err := filepath.Rel(filepath.Join(build.Default.GOPATH, "src"), directory)
	if err != nil {
		return nil, fmt.Errorf("provided directory not under GOPATH (%s): %v",
			build.Default.GOPATH, err)
	}

	conf := loader.Config{TypeChecker: types.Config{FakeImportC: true}}
	conf.Import(relDir)
	program, err := conf.Load()
	if err != nil {
		return nil, fmt.Errorf("couldn't load package: %v", err)
	}

	pkgInfo := program.Package(relDir)
	return &Package{
		Name:  pkgInfo.Pkg.Name(),
		files: pkgInfo.Files,
		defs:  pkgInfo.Defs,
	}, nil
}

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

// methodsOfTypeIn checks if a default methods is declared for the type,
// if declared - add it to the ignore list, and the template for this
// methods will NOT be added to the output file.
func (pkg *Package) methodsOfTypeIn(typeName string, decl *ast.FuncDecl) map[string]bool {
	if decl.Recv == nil || decl.Name == nil {
		return nil
	}

	var isTypeMethod, isPointerReceiver bool
	for _, field := range decl.Recv.List {
		if field.Type == nil {
			continue
		}

		var ok bool
		var ident *ast.Ident

		switch i := field.Type.(type) {
		case *ast.StarExpr:
			ident, ok = i.X.(*ast.Ident)
			isPointerReceiver = true
		case *ast.Ident:
			ident, ok = i, true
		}

		if !ok {
			continue
		}

		if ident.Name == typeName {
			isTypeMethod = true
		}

	}

	if !isTypeMethod {
		return nil
	}

	tmpls := map[string]bool{}
	for mName, shouldBePointer := range typeMethods {
		if !strings.Contains(decl.Name.Name, mName) {
			continue
		}

		if shouldBePointer == isPointerReceiver {
			tmpls[mName] = true
		}
	}

	return tmpls
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
