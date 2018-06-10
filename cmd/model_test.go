package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStructTags(t *testing.T) {
	res := parseStructTags(`db:"full_name" json:"full_name,omitempty"`)
	assert.Equal(t, "full_name", res["db"])
	assert.Equal(t, "full_name", res["json"])
}
