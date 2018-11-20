package templates

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTag(t *testing.T) {
	for _, value := range testDataStructTag {
		res, err := parseTag(value.tag)
		t.Run("check for tag=>"+value.tag, func(t *testing.T) {
			if value.throwError {
				if err == nil {
					fmt.Println(value.tag)
				}
				assert.NotNil(t, err, "error must be not nil")
				return
			}
			assert.Equal(t, value.result, res, "parsed tags mismatch")
		})
	}
}

var testDataStructTag = []struct {
	tag        string
	result     map[string]string
	throwError bool
}{
	{
		tag:        `db:"test"`,
		result:     map[string]string{"db": "test"},
		throwError: false,
	},
	{
		tag:        `db:"test" json:"testValue"`,
		result:     map[string]string{"db": "test", "json": "testValue"},
		throwError: false,
	},
	{
		tag:        `db:"test"`,
		throwError: false,
		result:     map[string]string{"db": "test"},
	},
	{
		tag:        `db:"test" json:"testValue" yaml:"value_42"`,
		throwError: false,
		result:     map[string]string{"db": "test", "json": "testValue", "yaml": "value_42"},
	},
	{
		tag:        "hello",
		result:     nil,
		throwError: true,
	}, // want "`hello` not compatible with reflect.StructTag.Get: bad syntax for struct tag pair"
	{
		tag:        "\tx:\"y\"",
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: bad syntax for struct tag key"
	{
		tag:        "x:\"y\"\tx:\"y\"",
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get"
	{
		tag:        "x:`y`",
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: bad syntax for struct tag value"
	{
		tag:        "ct\brl:\"char\"",
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: bad syntax for struct tag pair"
	{
		tag:        `:"emptykey"`,
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: bad syntax for struct tag key"
	{
		tag:        `x:"noEndQuote`,
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: bad syntax for struct tag value"
	{
		tag:        `x:"foo",y:"bar"`,
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: key:.value. pairs not separated by spaces"
	{
		tag:        `x:"foo"y:"bar"`,
		result:     nil,
		throwError: true,
	}, // want "not compatible with reflect.StructTag.Get: key:.value. pairs not separated by spaces"
	{
		tag:        `x:"y" u:"v" w:""`,
		result:     map[string]string{"x": "y", "u": "v", "w": ""},
		throwError: false,
	},
	{
		tag:        `x:"y:z" u:"v" w:""`,
		result:     map[string]string{"x": "y:z", "u": "v", "w": ""},
		throwError: false,
	}, // note multiple colons.
	{
		tag:        "k0:\"values contain spaces\" k1:\"literal\ttabs\" k2:\"and\\tescaped\\tabs\"",
		result:     map[string]string{"k0": "values contain spaces", "k1": "literal\ttabs", "k2": "and\\tescaped\\tabs"},
		throwError: false,
	},
	{
		tag:        `under_scores:"and" CAPS:"ARE_OK"`,
		result:     map[string]string{"under_scores": "and", "CAPS": "ARE_OK"},
		throwError: false,
	},
	{
		tag:        `json:"not_anon"`,
		result:     map[string]string{"json": "not_anon"},
		throwError: false,
	}, // ok; fields aren't embedded in JSON
	{
		tag:        `json:"-"`,
		result:     map[string]string{},
		throwError: false,
	}, // ok; entire field is ignored in JSON
	{
		tag:        `json:"a,omitempty"`,
		result:     map[string]string{"json": "a"},
		throwError: false,
	},
	{
		tag:        `json:"b, omitempty"`,
		result:     map[string]string{"json": "b"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `json:"c ,omitempty"`,
		result:     map[string]string{"json": "c "},
		throwError: false,
	},
	{
		tag:        `json:"d,omitempty, string"`,
		result:     map[string]string{"json": "d"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:"e local"`,
		result:     map[string]string{"xml": "e local"},
		throwError: false,
	},
	{
		tag:        `xml:"f "`,
		result:     map[string]string{"xml": "f "},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:" g"`,
		result:     map[string]string{"xml": " g"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:"h ,omitempty"`,
		result:     map[string]string{"xml": "h "},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:"i, omitempty"`,
		result:     map[string]string{"xml": "i"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:"j local ,omitempty"`,
		result:     map[string]string{"xml": "j local "},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:"k local, omitempty"`,
		result:     map[string]string{"xml": "k local"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:" l local,omitempty"`,
		result:     map[string]string{"xml": " l local"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:"m  local,omitempty"`,
		result:     map[string]string{"xml": "m  local"},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:" "`,
		result:     map[string]string{"xml": " "},
		throwError: false,
	}, // want "suspicious space in struct tag value"
	{
		tag:        `xml:""`,
		result:     map[string]string{"xml": ""},
		throwError: false,
	},
	{
		tag:        `xml:","`,
		result:     map[string]string{"xml": ""},
		throwError: false,
	},
	{
		tag:        `foo:" doesn't care "`,
		result:     map[string]string{"foo": " doesn't care "},
		throwError: false,
	},
	//{ // FIXME
	//	tag:        `x:"trunc\x0"`,
	//	result:     nil,
	//	throwError: true,
	//}, // want "not compatible with reflect.StructTag.Get: bad syntax for struct tag value"
}
