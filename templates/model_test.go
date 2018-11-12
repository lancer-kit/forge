package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTag(t *testing.T) {
	testData := []struct {
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
			result:     map[string]string{"db": "test"},
			throwError: false,
		},
		{
			tag:        `db:"test" json:"testValue" yaml:"value_42"`,
			result:     map[string]string{"db": "test", "json": "testValue", "yaml": "value_42"},
			throwError: false,
		},
	}
	for _, value := range testData {
		res, err := parseTag(value.tag)
		if value.throwError {
			assert.NotNil(t, err, "error must be not nil")
			continue
		}

		assert.Equal(t, value.result, res, "parsed tags mismatch")
	}
}
