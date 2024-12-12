package clapper

import (
	"strings"
	"unicode"
)

type TagType int

const (
	TagShort TagType = iota
	TagLong
	TagDefault
	TagHelp
	TagCommand
)

func GetTagType(tag string) (TagType, error) {
	switch tag {
	case "short":
		return TagShort, nil
	case "long":
		return TagLong, nil
	case "default":
		return TagDefault, nil
	case "help":
		return TagHelp, nil
	case "command":
		return TagCommand, nil
	default:
		return 0, NewUnknownTagTypeError(tag)
	}
}

// Tag is one property of a struct tag.
type Tag struct {
	// Type of the tag.
	Type TagType
	// Name gets derived from the struct field name if the tag is Short or Long and is pure computational.
	Name string
	// Value is an optional value given to the tag if an assignment operator is given. `short=s`
	Value string
	// Index of the tag found in the tag line.
	Index int
}

func NewTag(tag string, fieldName string, index int) (*Tag, error) {
	parts := strings.SplitN(tag, "=", 2)
	var value string
	if len(parts) == 2 {
		value = parts[1]
	}
	tagType, err := GetTagType(parts[0])
	if err != nil {
		return nil, err
	}
	result := &Tag{
		Type:  tagType,
		Value: value,
		Index: index,
	}

	if tagType == TagShort || tagType == TagLong {
		result.Name = result.DeriveName(fieldName)
	}

	return result, result.Validate()
}

func (t *Tag) validateShort() error {
	if len(t.Name) > 1 || len(t.Value) > 1 {
		return ErrShortOverrideCanOnlyBeOneLetter
	}
	return nil
}

func (t *Tag) validateLong() error {
	if len(t.Name) <= 1 || (t.HasValue() && len(t.Value) <= 1) {
		return ErrLongMustBeMoreThanOne
	}
	return nil
}

func (t *Tag) validateCommand() error {
	if len(t.Value) > 0 {
		return ErrCommandCanNotHaveValue
	}
	return nil
}

func (t *Tag) validateDefault() error {
	if len(t.Value) == 0 {
		return ErrNoDefaultValue
	}
	return nil
}

func (t *Tag) Validate() error {
	switch t.Type {
	case TagShort:
		return t.validateShort()
	case TagLong:
		return t.validateLong()
	case TagDefault:
		return t.validateDefault()
	case TagHelp:
		// No validation for help tags to not introduce breaking changes ATM.
	case TagCommand:
		return t.validateCommand()
	default:
		return NewUnknownTagTypeError(t.Name)
	}

	return nil
}

func (t *Tag) HasValue() bool {
	return t.Value != ""
}

func deriveLongName(fieldName string) string {
	var name string

	upperSequence := false
	sequenceCount := 0
	for i, char := range fieldName {
		lower := unicode.ToLower(char)

		newUpperSequence := unicode.IsUpper(char)
		if newUpperSequence == upperSequence {
			sequenceCount++
		}

		if newUpperSequence != upperSequence && i > 0 {
			if upperSequence && sequenceCount > 1 {
				l := len(name) - 1
				name = name[:l] + "-" + name[l:]
			}
			if !upperSequence {
				name += "-"
			}
			sequenceCount = 0
		}

		upperSequence = newUpperSequence

		name += string(lower)
	}

	return name
}

func (t *Tag) DeriveName(fieldName string) string {
	if t.HasValue() {
		if t.Type == TagShort {
			return t.Value[:1]
		}
		return t.Value
	}

	if t.Type == TagShort {
		return fieldName[:1]
	}

	name := deriveLongName(fieldName)

	return name
}
