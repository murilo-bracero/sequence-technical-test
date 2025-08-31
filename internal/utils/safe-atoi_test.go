package utils_test

import (
	"testing"

	"github.com/murilo-bracero/sequence-technical-test/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestSafeAtoi(t *testing.T) {
	table := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "valid input",
			input:    "123",
			expected: 123,
		},
		{
			name:     "invalid input",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "empty input",
			input:    "",
			expected: 0,
		},
	}

	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.SafeAtoi(tc.input, 0)
			assert.Equal(t, tc.expected, result)
		})
	}
}
