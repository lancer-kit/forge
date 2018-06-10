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

// JSONenums is a tool to automate the creation of methods that satisfy the
// fmt.Stringer, json.Marshaler and json.Unmarshaler interfaces.
// Given the name of a (signed or unsigned) integer type T that has constants
// defined, goplater will create a new self-contained Go source file implementing
//
//  func (t T) String() string
//  func (t T) MarshalJSON() ([]byte, error)
//  func (t *T) UnmarshalJSON([]byte) error
//
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// JSONenums is a simple implementation of a concept and the code might not be
// the most performant or beautiful to read.
//
// For example, given this snippet,
//
//	package painkiller
//
//	type Pill int
//
//	const (
//		Placebo Pill = iota
//		Aspirin
//		Ibuprofen
//		Paracetamol
//		Acetaminophen = Paracetamol
//	)
//
// running this command
//
//	goplater -type=Pill
//
// in the same directory will create the file pill_jsonenums.go, in package painkiller,
// containing a definition of
//
//  func (r Pill) String() string
//  func (r Pill) MarshalJSON() ([]byte, error)
//  func (r *Pill) UnmarshalJSON([]byte) error
//
// That method will translate the value of a Pill constant to the string representation
// of the respective constant name, so that the call fmt.Print(painkiller.Aspirin) will
// print the string "Aspirin".
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate goplater -type=Pill
//
// If multiple constants have the same value, the lexically first matching name will
// be used (in the example, Acetaminophen will print as "Paracetamol").
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is
// t_jsonenums.go, where t is the lower-cased name of the first type listed.
// The suffix can be overridden with the -suffix flag and a prefix may be added
// with the -prefix flag.
//
package main

import (
	"flag"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sheb-gregor/goplater/parser"
	"github.com/sheb-gregor/goplater/templates"
)

var (
	typeNames     = flag.String("type", "", "comma-separated list of type names; must be set")
	transformRule = flag.String("transform", "none", "enum item Name transformation method. Default: none")
	addTypePrefix = flag.Bool("tprefix", true, "add type name prefix into string values or not. Default: false")
	outputSuffix  = flag.String("suffix", "_enums", "suffix to be added to the output file")
	mergeSpecs    = flag.Bool("merge", false, "merge all output into one file. Default: false")
	//outputPrefix = flag.String("prefix", "", "prefix to be added to the output file")
	//forceRewritePrefix = flag.Bool("frw", false, "replace methods if they are implemented. Default: false")
)

func init() {
	flag.Parse()
	if len(*typeNames) == 0 {
		log.Fatalf("the flag -type must be set")
	}
	if transformRule == nil {
		log.Fatalf("transform method is not defined")
		return
	}
}

func main() {
	types := strings.Split(*typeNames, ",")

	// Only one directory at a time can be processed, and the default is ".".
	dir := "."
	if args := flag.Args(); len(args) == 1 {
		dir = args[0]
	} else if len(args) > 1 {
		log.Fatalf("only one directory at a time")
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("unable to determine absolute filepath for requested path %s: %v",
			dir, err)
	}

	if len(types) == 1 {
		*mergeSpecs = false
	}

	// need to remove already generated files for types
	// this is need for correct search of predefined by user
	// type vars and methods
	for _, typeName := range types {
		// Remove safe because we already check is path valid
		// and don't care about is present file - we need to remove it.
		os.Remove(getPath(typeName, dir))
	}

	if *mergeSpecs {
		os.Remove(getPath(mergeTypeNames(types), dir))
	}

	pkg, err := parser.ParsePackage(dir)
	if err != nil {
		log.Fatalf("parsing package: %v", err)
		return
	}

	var analysis = templates.Analysis{
		Command:     strings.Join(os.Args[1:], " "),
		PackageName: pkg.Name,
		Types:       make(map[string]templates.TypeSpec),
	}

	rule := TransformRule(*transformRule)

	// Run generate for each type.
	for _, typeName := range types {
		values, tmplsToExclude, err := pkg.ValuesOfType(typeName)
		if err != nil {
			log.Fatalf("finding values for type %v: %v", typeName, err)
		}
		analysis.Types[typeName] = templates.TypeSpec{
			TypeName:    typeName,
			Values:      rule.TransformValues(typeName, values, *addTypePrefix),
			ExcludeList: tmplsToExclude,
		}
	}

	for name, src := range analysis.GenerateByTemplate(*mergeSpecs) {
		if *mergeSpecs {
			name = mergeTypeNames(types)
		}

		if err := ioutil.WriteFile(getPath(name, dir), src, 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}

		if *mergeSpecs {
			return
		}
	}
}

func mergeTypeNames(names []string) string {
	sort.Strings(names)
	single := strings.Join(names, "_")
	crc32InUint32 := crc32.ChecksumIEEE([]byte(single))
	return strconv.FormatUint(uint64(crc32InUint32), 16)
}

func getPath(name, dir string) string {
	output := strings.ToLower(name + *outputSuffix + ".go")
	return filepath.Join(dir, output)
}
