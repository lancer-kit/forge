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

package templates

var FileBase = CodeTemplate{Name: "base", Raw: baseRaw}

var EnumBase = []CodeTemplate{
	{Name: "Invalid", Raw: typeError},
	{Name: "NameToValue", Raw: nameToValueRaw},
	{Name: "ValueToName", Raw: valueToNameRaw},
	{Name: "String", Raw: stringRaw},
	{Name: "Validate", Raw: validateRaw},
	{Name: "MarshalJSON", Raw: marshalJSONRaw},
	{Name: "UnmarshalJSON", Raw: unmarshalJSONRaw},
	{Name: "Value", Raw: rowValueRaw},
	{Name: "Scan", Raw: rowScanRaw},
}

var base = map[string]CodeTemplate{
	"base":          {Name: "base", Raw: baseRaw},
	"Invalid":       {Name: "Invalid", Raw: typeError},
	"NameToValue":   {Name: "NameToValue", Raw: nameToValueRaw},
	"ValueToName":   {Name: "ValueToName", Raw: valueToNameRaw},
	"String":        {Name: "String", Raw: stringRaw},
	"Validate":      {Name: "Validate", Raw: validateRaw},
	"MarshalJSON":   {Name: "MarshalJSON", Raw: marshalJSONRaw},
	"UnmarshalJSON": {Name: "UnmarshalJSON", Raw: unmarshalJSONRaw},
	"Value":         {Name: "Value", Raw: rowValueRaw},
	"Scan":          {Name: "Scan", Raw: rowScanRaw},
}

//
//func init() {
//	FileBase.parse()
//
//	for i := range EnumBase {
//		EnumBase[i].parse()
//	}
//}

var (
	baseRaw = `
// generated by goplater {{.Command}}; DO NOT EDIT
package {{.PackageName}}

import (
    "database/sql"
    "database/sql/driver"
    "encoding/json"
    "errors"
    "fmt"
)

func init() {
    // stub usage of json for situation when
    // (Un)MarshalJSON methods will be omitted
    _ = json.Delim('s')

    // stub usage of sql/driver for situation when
    // Scan/Value methods will be omitted
    _ = driver.Bool
    _ = sql.LevelDefault
}
`
	typeError = `
var Err{{.TypeName}}Invalid = errors.New("{{.TypeName}} is invalid")
`
	nameToValueRaw = `
var def{{.TypeName}}NameToValue = map[string]{{.TypeName}} {
    {{range .Values}}"{{.Str}}": {{.Name}},
    {{end}}
}
`

	valueToNameRaw = `
var def{{.TypeName}}ValueToName = map[{{.TypeName}}]string {
        {{range .Values}}{{.Name}}: "{{.Str}}",
        {{end}}
    }
`
	stringRaw = `
// String is generated so {{.TypeName}} satisfies fmt.Stringer.
func (r {{.TypeName}}) String() string {
    s, ok := def{{.TypeName}}ValueToName[r]
    if !ok {
        return fmt.Sprintf("{{.TypeName}}(%d)", r)
    }
    return s
}
`

	validateRaw = `
// Validate verifies that value is predefined for {{.TypeName}}.
func (r {{.TypeName}}) Validate() error {
    _, ok := def{{.TypeName}}ValueToName[r]
    if !ok {
        return Err{{.TypeName}}Invalid
    }
    return nil
}
`

	marshalJSONRaw = `
// MarshalJSON is generated so {{.TypeName}} satisfies json.Marshaler.
func (r {{.TypeName}}) MarshalJSON() ([]byte, error) {
    if s, ok := interface{}(r).(fmt.Stringer); ok {
        return json.Marshal(s.String())
    }
    s, ok := def{{.TypeName}}ValueToName[r]
    if !ok {
        return nil, fmt.Errorf("{{.TypeName}}(%d) is invalid value", r)
    }
    return json.Marshal(s)
}
`

	unmarshalJSONRaw = `
// UnmarshalJSON is generated so {{.TypeName}} satisfies json.Unmarshaler.
func (r *{{.TypeName}}) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return fmt.Errorf("{{.TypeName}}: should be a string, got %s", string(data))
    }
    v, ok := def{{.TypeName}}NameToValue[s]
    if !ok {
        return fmt.Errorf("{{.TypeName}}(%q) is invalid value", s)
    }
    *r = v
    return nil
}
`

	rowValueRaw = `
// Value is generated so {{.TypeName}} satisfies db row driver.Valuer.
func (r {{.TypeName}}) Value() (driver.Value, error) {
    s, ok := def{{.TypeName}}ValueToName[r]
    if !ok {
        return nil, nil
    }
    return s, nil
}
`

	rowScanRaw = `
// Value is generated so {{.TypeName}} satisfies db row driver.Scanner.
func (r *{{.TypeName}}) Scan(src interface{}) error {
    switch v := src.(type) {
    case string:
        val, _ := def{{.TypeName}}NameToValue[v]
        *r = val
        return nil
    case []byte:
        var i {{.TypeName}}
        err := json.Unmarshal(v, &i)
        if err != nil {
            return errors.New("{{.TypeName}}: can't unmarshal column data")
        }
    
        *r = i
        return nil
    case int, int8, int32, int64, uint, uint8, uint32, uint64:
        ni := sql.NullInt64{}
        err := ni.Scan(v)
        if err != nil {
            return errors.New("{{.TypeName}}: can't scan column data into int64")
        }
    
        *r = {{.TypeName}}(ni.Int64)
        return nil
    }
    return errors.New("{{.TypeName}}: invalid type")
}
`
)
