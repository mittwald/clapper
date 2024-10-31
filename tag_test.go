package clapper

import "testing"

func TestDeriveLongName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "single letter", input: "a", expected: "a"},
		{name: "multiple letters", input: "abc", expected: "abc"},
		{name: "mixed case", input: "AbcDeF", expected: "abc-de-f"},
		{name: "empty string", input: "", expected: ""},
		{name: "single underscore", input: "_", expected: "_"},
		{name: "single hyphen", input: "-", expected: "-"},
		{name: "some exported", input: "FooBar", expected: "foo-bar"},
		{name: "acronym at start", input: "UDPMode", expected: "udp-mode"},
		{name: "acronym at end", input: "someAPI", expected: "some-api"},
		{name: "acronym inside", input: "ExternalTCPSocket", expected: "external-tcp-socket"},
		{name: "lower start", input: "fooBar", expected: "foo-bar"},
		{name: "some numbers", input: "123hello", expected: "123hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deriveLongName(tt.input); got != tt.expected {
				t.Errorf("deriveLongName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
