package main

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestTransformRule_Transform(t *testing.T) {
	testCases := []struct {
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

	for _, tCase := range testCases {
		for rule, result := range tCase.expected {
			assert.Equal(t, result, rule.Transform(tCase.value),
				fmt.Sprintf("error in case for rule (%s) with value (%s)", rule, tCase.value))
		}
	}
}
