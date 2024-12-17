package clapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArgParserExtGet(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		key     string
		argType ArgType
		wantOk  bool
		want    []string
	}{
		{
			name:    "single short flag without value",
			input:   []string{"-d"},
			key:     "d",
			argType: ArgTypeShort,
			wantOk:  true,
			want:    []string{},
		},
		{
			name:    "single short flag with value",
			input:   []string{"-d", "hello"},
			key:     "d",
			argType: ArgTypeShort,
			wantOk:  true,
			want:    []string{"hello"},
		},
		{
			name:    "single short flag with values",
			input:   []string{"-d", "hello", "world", "--other"},
			key:     "d",
			argType: ArgTypeShort,
			wantOk:  true,
			want:    []string{"hello", "world"},
		},
		{
			name:    "repeated short flag with value",
			input:   []string{"-d", "hello", "--other", "-d", "world"},
			key:     "d",
			argType: ArgTypeShort,
			wantOk:  true,
			want:    []string{"hello", "world"},
		},
		{
			name:    "not existing",
			input:   []string{"-d", "hello", "--foo"},
			key:     "bar",
			argType: ArgTypeLong,
			wantOk:  false,
			want:    []string{},
		},
		{
			name:    "existing long",
			input:   []string{"-d", "hello", "--foo", "bar"},
			key:     "foo",
			argType: ArgTypeLong,
			wantOk:  true,
			want:    []string{"bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewArgParserExt(tt.input)
			got, ok := parser.Get(tt.key, tt.argType)
			require.Equal(t, tt.wantOk, ok, "got unexpected Get ok result")
			assert.ElementsMatch(t, tt.want, got, "got unexpected Get values")
		})
	}
}

func TestTrailing_SingleEndValue(t *testing.T) {
	parser := NewArgParserExt([]string{"-d", "hello", "world", "--other", "foo", "bar", "baz"})
	got, ok := parser.Get("other", ArgTypeLong)
	require.True(t, ok, "got unexpected Get ok result")
	assert.Equal(t, []string{"foo", "bar", "baz"}, got)
	parser.Consume("other", ArgTypeLong, 1)
	assert.Equal(t, []string{"bar", "baz"}, parser.GetTrailing())
}

func TestTrailing_SliceEndValue(t *testing.T) {
	parser := NewArgParserExt([]string{"-d", "hello", "world", "--other", "foo", "bar", "baz"})
	got, ok := parser.Get("other", ArgTypeLong)
	require.True(t, ok, "got unexpected Get ok result")
	assert.Equal(t, []string{"foo", "bar", "baz"}, got)
	parser.Consume("other", ArgTypeLong, len(got))
	assert.Equal(t, []string{}, parser.GetTrailing())
}

func TestConsumeTrailing(t *testing.T) {
	parser := NewArgParserExt([]string{"-d", "hello", "world", "--other", "foo", "bar", "baz"})
	parser.ConsumeTrailing(1)
	assert.Equal(t, []string{"bar", "baz"}, parser.GetTrailing())
}

func TestConsumeAllTrailing(t *testing.T) {
	parser := NewArgParserExt([]string{"-d", "hello", "world", "--other", "foo", "bar", "baz"})
	parser.ConsumeTrailing(3)
	assert.Equal(t, []string{}, parser.GetTrailing())
}
