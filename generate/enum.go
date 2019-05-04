package generate

import (
	"flag"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/lancer-kit/forge/configs"
	"github.com/lancer-kit/forge/parser"
	"github.com/lancer-kit/forge/templates"
)

func Enums(config configs.EnumsConfig) error {
	// Only one directory at a time can be processed, and the default is ".".
	dir := "."
	if args := flag.Args(); len(args) >= 1 {
		dir = args[0]
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("unable to determine absolute filepath for requested path %s: %v", dir, err)
	}

	if len(config.Types) == 1 {
		config.MergeSpecs = false
	}

	// need to remove already generated files for types
	// this is need for correct search of predefined by user
	// type vars and methods
	for _, typeName := range config.Types {
		// Remove safe because we already check is path valid
		// and don't care about is present file - we need to remove it.
		os.Remove(config.GetPath(typeName, dir))
	}

	if config.MergeSpecs {
		os.Remove(config.GetPath(mergeTypeNames(config.Types), dir))
	}

	pkg, err := parser.ParsePackage(dir)
	if err != nil {
		return fmt.Errorf("parsing package: %v", err)
	}

	var analysis = templates.Analysis{
		Command:     strings.Join(os.Args[1:], " "),
		PackageName: pkg.Name,
		Types:       make(map[string]templates.TypeSpec),
	}

	rule := templates.TransformRule(config.TransformRule)

	// Run generate for each type.
	for _, typeName := range config.Types {
		values, tmplsToExclude, err := pkg.ValuesOfType(typeName)
		if err != nil {
			return fmt.Errorf("finding values for type %v: %v", typeName, err)
		}
		analysis.Types[typeName] = templates.TypeSpec{
			TypeName:    typeName,
			Values:      rule.TransformValues(typeName, values, config.AddTypePrefix),
			ExcludeList: tmplsToExclude,
		}
	}

	for name, src := range analysis.GenerateByTemplate(config.MergeSpecs) {
		if config.MergeSpecs {
			name = mergeTypeNames(config.Types)
		}

		if err := ioutil.WriteFile(config.GetPath(name, dir), src, 0644); err != nil {
			return fmt.Errorf("writing output: %s", err)
		}

		if config.MergeSpecs {
			return nil
		}
	}

	return nil
}

func mergeTypeNames(names []string) string {
	sort.Strings(names)
	single := strings.Join(names, "_")
	crc32InUint32 := crc32.ChecksumIEEE([]byte(single))
	return strconv.FormatUint(uint64(crc32InUint32), 16)
}
