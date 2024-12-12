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

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		tag     Tag
		wantErr bool
	}{
		{name: "long tag with long value is ok", tag: Tag{Type: TagLong, Name: "foo", Value: "foo"}, wantErr: false},
		{name: "long tag with short value fails", tag: Tag{Type: TagLong, Name: "foo", Value: "f"}, wantErr: true},
		{name: "short tag with long value fails", tag: Tag{Type: TagShort, Name: "f", Value: "foo"}, wantErr: true},
		{name: "short tag with short value is ok", tag: Tag{Type: TagShort, Name: "f", Value: "f"}, wantErr: false},
		{name: "command tag with value fails", tag: Tag{Type: TagCommand, Name: "", Value: "some"}, wantErr: true},
		{name: "command tag without value is ok", tag: Tag{Type: TagCommand, Name: "", Value: ""}, wantErr: false},
		{name: "default tag with value is ok", tag: Tag{Type: TagDefault, Name: "", Value: "some"}, wantErr: false},
		{name: "default tag without value fails", tag: Tag{Type: TagDefault, Name: "", Value: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tag.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagType(t *testing.T) {
	tests := []struct {
		tagName     string
		wantTagType TagType
		wantErr     bool
	}{
		{tagName: "short", wantTagType: TagShort, wantErr: false},
		{tagName: "long", wantTagType: TagLong, wantErr: false},
		{tagName: "default", wantTagType: TagDefault, wantErr: false},
		{tagName: "help", wantTagType: TagHelp, wantErr: false},
		{tagName: "command", wantTagType: TagCommand, wantErr: false},
		{tagName: "unknown", wantTagType: 0, wantErr: true},
		{tagName: "SHORT", wantTagType: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.tagName, func(t *testing.T) {
			gotTagType, err := GetTagType(tt.tagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTagType != tt.wantTagType {
				t.Errorf("GetTagType() = %v, want %v", gotTagType, tt.wantTagType)
			}
		})
	}
}
