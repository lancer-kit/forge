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

```text
$  ./forge 
NAME:
   forge - cli tool and generator from lancer-kit

USAGE:
   forge [global options] command [command options] [arguments...]

VERSION:
   2.5

COMMANDS:
     enum     generate var and methods for the iota-enums
     model    generate code for structure by template
     bindata  forge bindata <options>
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

### Bindata

Build-in fork of [go-bindata](https://github.com/jteeuwen/go-bindata)


```text
USAGE:
   forge bindata [command options] [arguments...]


OPTIONS:
   --debug           Do not embed the assets, but provide the embedding API. Contents will still be loaded from disk.
   --dev             Similar to debug, but does not emit absolute paths. Expects a rootDir variable to already exist in the generated code's package.
   --nomemcopy       Use a .rodata hack to get rid of unnecessary memcopies. Refer to the documentation to see what implications this carries.
   --nocompress      Assets will *not* be GZIP compressed when this flag is specified.
   --nometadata      Assets will not preserve size, mode, and modtime info.
   --tags value      Optional set of build tags to include.
   --prefix value    Optional path prefix to strip off asset names.
   --pkg value       Package name to use in the generated code.
   -o value          Optional name of the output file to be generated.
   --mode value      Optional file mode override for all files. (default: 0)
   --modetime value  Optional modification unix timestamp override for all files. (default: 0)
   --ignore value    Regex pattern to ignore
   -i value          List of input directories/files
```


# ToDo

- [ ] Improve documentation, add tests
- [ ] Add string types support 
- [ ] Add bitmap types support
- [x] Add custom template support


# License

This tool contains code from next repos: 
  - https://github.com/campoy/jsonenums.
  - https://github.com/jteeuwen/go-bindata.

