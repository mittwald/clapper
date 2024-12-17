package clapper

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	_ error = ParseError{}
	_ error = UnsupportedReflectTypeError{}
	_ error = UnknownTagTypeError{}
	_ error = UnexpectedInputFormatError{}
	_ error = CommandRequiredError{}

	ErrNoStruct                        = errors.New("target is not a struct")
	ErrEmptyArgument                   = errors.New("empty argument")
	ErrUnexpectedValue                 = errors.New("expected flag not value")
	ErrFieldCanNotBeSet                = errors.New("field can't be set")
	ErrNoFlagSpecifier                 = errors.New("no flag specified for struct field")
	ErrShortOverrideCanOnlyBeOneLetter = errors.New("short override can only be one letter")
	ErrLongMustBeMoreThanOne           = errors.New("long name must be more than one character")
	ErrCommandCanNotHaveValue          = errors.New("command can't have a value")
	ErrDuplicateCommandTag             = errors.New("duplicate command tag found")
	ErrNoDefaultValue                  = errors.New("default spcified but no default value given")

	errInternalNoArgumentsForTag = errors.New("no arguments for tag found")
)

// CommandRequiredError will be thrown when a command tag is required but no command is provided.
// A `help` given in the tag will enrich the error message with the given help message.
type CommandRequiredError struct {
	help string
}

func NewCommandRequiredError(help string) CommandRequiredError {
	return CommandRequiredError{help: help}
}

func (e CommandRequiredError) Error() string {
	result := "command required but no command given"
	if len(e.help) > 0 {
		result += fmt.Sprintf(". Use: %s", e.help)
	}

	return result
}

// UnsupportedReflectTypeError will be thrown when a struct field has a type that can not be set with the provided value.
// For example givving a string to a field of type int.
type UnexpectedInputFormatError struct {
	Input          string
	ExpectedFormat reflect.Type
}

func NewUnexpectedInputFormatError(input string, expected reflect.Type) UnexpectedInputFormatError {
	return UnexpectedInputFormatError{
		Input:          input,
		ExpectedFormat: expected,
	}
}

func (e UnexpectedInputFormatError) Error() string {
	return fmt.Sprintf("unexpected input format. given '%s', expected %s", e.Input, e.ExpectedFormat)
}

// UnsupportedReflectTypeError will be thrown when a struct tag given is unknown to claper.
type UnknownTagTypeError struct {
	TagName string
}

func NewUnknownTagTypeError(tagName string) UnknownTagTypeError {
	return UnknownTagTypeError{TagName: tagName}
}

func (e UnknownTagTypeError) Error() string {
	return fmt.Sprintf("unknown struct tag '%s'", e.TagName)
}

// MandatoryParameterError will be thrown when a mandatory parameter is missing and no default value is provided.
type MandatoryParameterError struct {
	Name string
}

func (e MandatoryParameterError) Error() string {
	return fmt.Sprintf("required parameter '%s' is missing and no default available", e.Name)
}

func NewMandatoryParameterError(name string) MandatoryParameterError {
	return MandatoryParameterError{Name: name}
}

// ParseError will be thrown when an error occurs during parsing.
type ParseError struct {
	error
	Index   int
	Name    string
	TagLine string
}

func NewParseError(from error, index int, name string, tagLine string) ParseError {
	return ParseError{
		error:   from,
		Index:   index,
		Name:    name,
		TagLine: tagLine,
	}
}

func (e ParseError) Underlying() error {
	return e.error
}

func (e ParseError) Error() string {
	return fmt.Sprintf("parse error '%s' at index %d: field '%s' (tag-line: %s)", e.error, e.Index, e.Name, e.TagLine)
}

// UnknownTagTypeError will be thrown when a struct field type is unsupported by claper.
type UnsupportedReflectTypeError struct {
	Type string
}

func NewUnsupportedReflectTypeError(t string) UnsupportedReflectTypeError {
	return UnsupportedReflectTypeError{Type: t}
}

func (e UnsupportedReflectTypeError) Error() string {
	return "unsupported reflect type:" + e.Type
}
