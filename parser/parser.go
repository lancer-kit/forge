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
	"go/types"
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
