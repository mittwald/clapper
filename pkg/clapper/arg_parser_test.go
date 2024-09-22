package clapper

import (
	"testing"
)

func TestArguments(t *testing.T) {
	test := []struct {
		name     string
		args     []string
		expected *ArgsParser
	}{
		{
			name:     "no arguments at all",
			args:     []string{},
			expected: NewArgsParser(),
		},
		{
			name: "no flags only values",
			args: []string{"foo", "bar"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Trailing = []string{"foo", "bar"}
			}),
		},
		{
			name: "combined shorts",
			args: []string{"-fun"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{}
				args.Params["u"] = []string{}
				args.Params["n"] = []string{}
			}),
		},
		{
			name: "multiple shorts without value",
			args: []string{"-f", "-u", "-n"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{}
				args.Params["u"] = []string{}
				args.Params["n"] = []string{}
			}),
		},
		{
			name: "repeated shorts without value",
			args: []string{"-f", "-f", "-f"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{}
			}),
		},
		{
			name: "multiple and combined shorts without value",
			args: []string{"-f", "-fllloool", "-f"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{}
				args.Params["l"] = []string{}
				args.Params["o"] = []string{}
			}),
		},
		{
			name: "combined repeated shorts with value binds to last flag",
			args: []string{"-flol", "123"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{}
				args.Params["l"] = []string{"123"}
				args.Params["o"] = []string{}
				args.Trailing = []string{"123"}
			}),
		},
		{
			name: "repeated shorts with value",
			args: []string{"-f", "1", "-f", "2", "-f", "42"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{"1", "2", "42"}
				args.Trailing = []string{"42"}
			}),
		},
		{
			name: "repeated shorts without value and extra value",
			args: []string{"-f", "1", "-f", "2", "-f", "42", "extra", "value"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["f"] = []string{"1", "2", "42", "extra", "value"}
				args.Trailing = []string{"42", "extra", "value"}
			}),
		},
		{
			name: "repeated long without value",
			args: []string{"--foo", "--foo"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["foo"] = []string{}
			}),
		},
		{
			name: "multiple long without value",
			args: []string{"--foo", "--bar"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["foo"] = []string{}
				args.Params["bar"] = []string{}
			}),
		},
		{
			name: "repeated long with value",
			args: []string{"--foo", "bar", "--foo", "baz"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["foo"] = []string{"bar", "baz"}
				args.Trailing = []string{"baz"}
			}),
		},
		{
			name: "long with multiple values",
			args: []string{"--foo", "bar", "baz"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["foo"] = []string{"bar", "baz"}
				args.Trailing = []string{"bar", "baz"}
			}),
		},
		{
			name: "arguments with trailing flags",
			args: []string{"hello", "world", "--foo"},
			expected: NewArgsParser().With(func(args *ArgsParser) {
				args.Params["foo"] = []string{}
			}),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewArgsParser()

			_, err := parser.Parse(tt.args)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !parser.ValuesEqualWith(tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, parser)
			}
		})
	}
}
