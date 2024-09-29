package clapper

import (
	"errors"
	"fmt"
)

var (
	_ error = ParseError{}
	_ error = UnsupportedReflectTypeError{}

	ErrNoStruct                        = errors.New("target is not a struct")
	ErrEmptyArgument                   = errors.New("empty argument")
	ErrUnexpectedValue                 = errors.New("expected flag not value")
	ErrFieldCanNotBeSet                = errors.New("field can't be set")
	ErrNoFlagSpecifier                 = errors.New("no flag specified for struct field")
	ErrShortOverrideCanOnlyBeOneLetter = errors.New("short override can only be one letter")
	ErrLongMustBeMoreThanOne           = errors.New("long name must be more than one character")
)

type MandatoryParameterError struct {
	Name string
}

func (e MandatoryParameterError) Error() string {
	return fmt.Sprintf("required parameter '%s' is missing and no default available", e.Name)
}

func NewMandatoryParameterError(name string) MandatoryParameterError {
	return MandatoryParameterError{Name: name}
}

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

type UnsupportedReflectTypeError struct {
	Type string
}

func NewUnsupportedReflectTypeError(t string) UnsupportedReflectTypeError {
	return UnsupportedReflectTypeError{Type: t}
}

func (e UnsupportedReflectTypeError) Error() string {
	return "unsupported reflect type:" + e.Type
}
