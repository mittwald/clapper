package clapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSanitizeArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no arguments",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "leading values",
			input:    []string{"foo", "bar", "-d"},
			expected: []string{"-d"},
		},
		{
			name:     "combined short flags",
			input:    []string{"-d", "-abc"},
			expected: []string{"-d", "-a", "-b", "-c"},
		},
		{
			name:     "args with values",
			input:    []string{"-d", "hello", "--some"},
			expected: []string{"-d", "hello", "--some"},
		},
		{
			name:     "args with assigned values",
			input:    []string{"-d=hello", "--some"},
			expected: []string{"-d", "hello", "--some"},
		},
		{
			name:     "trailing values",
			input:    []string{"-d", "foo", "bar"},
			expected: []string{"-d", "foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultArgumentSanitizer(tt.input).Get()
			assert.ElementsMatch(t, tt.expected, got, "got unexpected sanitizeArgs result")
		})
	}
}
