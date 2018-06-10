# goplater

**goplater** is a tool for automating the creation of methods that satisfy some built-in interfaces or those defined in the templates.

# Usage 

**goplater** can be run from console in directory with target type or with help `go generate`:

```bash
$ goplater enum --type MyType
``` 

Or

``` go
package main

//go:generate goplater enum --type MyType
type MyType int 

const (
	MyTypeA MyType   = iota
	MyTypeB
	MyTypeC
)
```

List of arguments:

| Flag | Values | Description |
| ---- | ------ | ----------- |
| type | string | The name of the target type or types for code generation |
| transform | snake, kebab, space, none | A rule describing the strategy for converting constant names to a string. Default: none|
| tprefix | true, false | add type name prefix into string values or not. Default: false |
| prefix | string |  A prefix to be added to the output file |
| suffix | string |  A suffix to be added to the output. Default: "_enums"|

# ToDo

- [ ] Improve documentation, add tests
- [ ] Add string types support 
- [ ] Add bitmap types support
- [ ] Add custom template support


# License

This package is based on https://github.com/campoy/jsonenums.