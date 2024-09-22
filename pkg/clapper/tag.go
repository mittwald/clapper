package clapper

import (
	"fmt"
	"strings"
)

type TagType int

const (
	TagShort TagType = iota
	TagLong
	TagDefault
	TagHelp
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
	default:
		return 0, fmt.Errorf("unknown tag: %s", tag)
	}
}

type Tag struct {
	Type  TagType
	Name  string
	Value string
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

func (t *Tag) Validate() error {
	if t.Type == TagLong {
		if len(t.Name) <= 1 || (t.HasValue() && len(t.Value) <= 1) {
			return ErrLongMustBeMoreThanOne
		}
	}

	if t.Type == TagShort {
		if len(t.Name) > 1 || len(t.Value) > 1 {
			return ErrShortOverrideCanOnlyBeOneLetter
		}
	}

	return nil
}

func (t *Tag) HasValue() bool {
	return t.Value != ""
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

	var name string
	for i, char := range fieldName {
		lower := rune(strings.ToLower(string(char))[0])
		if char != lower && i != 0 && i != len(fieldName)-1 {
			name += "-"
		}
		name += string(lower)
	}

	return name
}
