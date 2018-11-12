package templates

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	value    string
	expected map[TransformRule]string
}{
	{
		value: "Test",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test",
			TransformRuleKebab:       "test",
			TransformRuleSpace:       "test",
			TransformRuleNone:        "Test",
			TransformRule("invalid"): "Test",
		},
	},
	{
		value: "TestValue",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test_value",
			TransformRuleKebab:       "test-value",
			TransformRuleSpace:       "test value",
			TransformRuleNone:        "TestValue",
			TransformRule("invalid"): "TestValue",
		},
	},
	{
		value: "TestQ",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test_q",
			TransformRuleKebab:       "test-q",
			TransformRuleSpace:       "test q",
			TransformRuleNone:        "TestQ",
			TransformRule("invalid"): "TestQ",
		},
	},
	{
		value: "TestValueQ",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test_value_q",
			TransformRuleKebab:       "test-value-q",
			TransformRuleSpace:       "test value q",
			TransformRuleNone:        "TestValueQ",
			TransformRule("invalid"): "TestValueQ",
		},
	},
	{
		value: "TestQValue",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test_q_value",
			TransformRuleKebab:       "test-q-value",
			TransformRuleSpace:       "test q value",
			TransformRuleNone:        "TestQValue",
			TransformRule("invalid"): "TestQValue",
		},
	},
	{
		value: "Testvalueq",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "testvalueq",
			TransformRuleKebab:       "testvalueq",
			TransformRuleSpace:       "testvalueq",
			TransformRuleNone:        "Testvalueq",
			TransformRule("invalid"): "Testvalueq",
		},
	},
	{
		value: "testValueq",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test_valueq",
			TransformRuleKebab:       "test-valueq",
			TransformRuleSpace:       "test valueq",
			TransformRuleNone:        "testValueq",
			TransformRule("invalid"): "testValueq",
		},
	},
	{
		value: "TTL",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "ttl",
			TransformRuleKebab:       "ttl",
			TransformRuleSpace:       "ttl",
			TransformRuleNone:        "TTL",
			TransformRule("invalid"): "TTL",
		},
	},
	{
		value: "TestTTL",
		expected: map[TransformRule]string{
			TransformRuleSnake:       "test_ttl",
			TransformRuleKebab:       "test-ttl",
			TransformRuleSpace:       "test ttl",
			TransformRuleNone:        "TestTTL",
			TransformRule("invalid"): "TestTTL",
		}},
}

func TestTransformRule_Transform(t *testing.T) {
	for _, tCase := range testCases {
		for rule, result := range tCase.expected {
			assert.Equal(t, result, rule.Transform(tCase.value),
				fmt.Sprintf("error in case for rule (%s) with value (%s)", rule, tCase.value))
		}
	}
}

func TestTransformRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    TransformRule
		wantErr bool
	}{
		{name: "valid:snake", rule: TransformRuleSnake, wantErr: false},
		{name: "valid:kebab", rule: TransformRuleKebab, wantErr: false},
		{name: "valid:space", rule: TransformRuleSpace, wantErr: false},
		{name: "valid:none", rule: TransformRuleNone, wantErr: false},
		{name: "invalid", rule: "test_value", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.rule.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TransformRule.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
