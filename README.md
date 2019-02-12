# forge

//todo: update description
**forge** is a tool for automating the creation of methods that satisfy some built-in interfaces or those defined in the templates.

# Usage 

**forge** can be run from console in directory with target type or with help `go generate`:

```bash
$ forge enum --type MyType
``` 

Or

``` go
package main

//go:generate forge enum --type MyType
type MyType int 

const (
	MyTypeA MyType   = iota
	MyTypeB
	MyTypeC
)
```
```bash
$  ./forge 
NAME:
   forge - auto generate, dont repeat

USAGE:
   forge [global options] command [command options] [arguments...]

VERSION:
   2.0

COMMANDS:
     enum     generate var and methods for the iota-enums
     model    generate code for structure by template
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

```

### Enum

Command: `forge enum`

Implements predefined interfaces for a user-defined integer type with **iota**-constants.

Interface list:
  
1. fmt.Stringer - `String() string`;
2. json.Marshaler - `MarshalJSON() ([]byte, error)`;
3. json.Unmarshaler - `UnmarshalJSON([]byte) error`;
4. driver.Valuer - `Value() (Value, error)`;
5. sql.Scanner - `Scan(src interface{}) error`;
6. Validator - `Validate() error`.

Predefined variables :

- `var def<Type>ValueToName map[<Type>]string` - matching a constant and its string representation;
- `var def<Type>NameToValue map[string]<Type>` - matching a string representation and constant; 
- `var Err<Type>Invalid error` - error.

All methods and maps can be pre-determined before generation, and at run they will be omitted.

List of arguments:

| Flag | Type | Description |
| ---- | ------ | ----------- |
| type | string | The name of the target type or types for code generation |
| transform | snake, kebab, space, none | A rule describing the strategy for converting constant names to a string. Default: none|
| tprefix | true, false | add type name prefix into string values or not. Default: false |
| prefix | string |  A prefix to be added to the output file |
| suffix | string |  A suffix to be added to the output. Default: "_enums"|
| merge | bool |  Merge all output into one file, if set `prefix` and `suffix` will be ignored. Default: false|

Example
```bash
forge enum --type ShirtSize,WeekDay --merge true
```

### Model 

Command: `forge model`

Find the structure definition, analyze it and fill in the transposed template.


Available fields for template:

| Flag | Type | Description |
| ---- | ------ | ----------- |
| Package | string  | Name of package with type |
| TypeName | string   | Type  name|
| TypeString | string |  camelCased type name | 
| Fields | []`Field` |  List of structure field definitions |
| `Field`.Name | string | Name of filed |
| `Field`.FType | string | Type of field |
| `Field`.Tags | map[string]string | Field tags |


List of arguments:

| Flag | Type | Description |
| ---- | ------ | ----------- |
| tmpl | string  |   path to the templates; required |
| type | string   |  list of type names; required |
| prefix | string |  prefix to be added to the output file | 
| suffix | string |  suffix to be added to the output file | 


# ToDo

- [ ] Improve documentation, add tests
- [ ] Add string types support 
- [ ] Add bitmap types support
- [ ] Add custom template support


# License

This tool contains code from next repos: 
  - https://github.com/campoy/jsonenums.
  - https://github.com/jteeuwen/go-bindata.

