package clapper

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	TagName = "clapper"
)

type ParsedTags = map[int]map[TagType]Tag

func parseTags(tagItems []string, fieldName string, index int) (map[TagType]Tag, error) {
	tags := make(map[TagType]Tag, 0)
	for _, tagItem := range tagItems {
		tag, err := NewTag(tagItem, fieldName, index)
		if err != nil {
			return nil, err
		}
		tags[tag.Type] = *tag
	}
	return tags, nil
}

func parseStructTags(t reflect.Type) (ParsedTags, error) {
	parsedTags := make(map[int]map[TagType]Tag, 0)
	commandTagSpecified := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagLine := field.Tag.Get(TagName)
		if tagLine == "" {
			continue
		}
		tagItems := strings.Split(tagLine, ",")
		tags, err := parseTags(tagItems, field.Name, i)
		if err != nil {
			return nil, NewParseError(err, i, field.Name, tagLine)
		}
		if hasTagType(tags, TagCommand) {
			if commandTagSpecified {
				return nil, ErrDuplicateCommandTag
			}
			commandTagSpecified = true
		}
		parsedTags[i] = tags
	}
	return parsedTags, nil
}

func paramName(tags map[TagType]Tag) string {
	tag, ok := tags[TagLong]
	if !ok {
		tag, ok = tags[TagShort]
		if !ok {
			return "<unknown>"
		}
	}
	return tag.Name
}

func mustTagTypeToArgType(tagType TagType) ArgType {
	switch tagType {
	case TagShort:
		return ArgTypeShort
	case TagLong:
		return ArgTypeLong
	default:
		panic(fmt.Sprintf("unknown tag type to map to arg type: %v", tagType))
	}
}

func valuesFor(tagType TagType, tags map[TagType]Tag, args *ArgParserExt) (key string, values []string) {
	tag, ok := tags[tagType]
	if !ok {
		return "", nil
	}
	argType := mustTagTypeToArgType(tag.Type)

	key = tag.GetLookupKey()

	values, ok = args.Get(key, argType)
	if !ok {
		return key, nil
	}
	return key, values
}

func isPointer(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Ptr
}

func isBool(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Bool
}

func isOptionalField(field reflect.StructField) bool {
	return isPointer(field) || isBool(field)
}

func trySetForType(tagType TagType, field reflect.StructField, fieldValue reflect.Value, tags map[TagType]Tag, args *ArgParserExt) error {
	key, values := valuesFor(tagType, tags, args)
	if values == nil {
		return errInternalNoArgumentsForTag
	}

	took, err := StringReflect(field, fieldValue, values)
	if err != nil {
		return err
	}

	argType := mustTagTypeToArgType(tagType)
	args.Consume(key, argType, took)

	return nil
}

func trySetDefault(field reflect.StructField, fieldValue reflect.Value, tags map[TagType]Tag) error {
	tag, ok := tags[TagDefault]
	if !ok {
		if isOptionalField(field) {
			return nil
		}
		return NewMandatoryParameterError(paramName(tags))
	}
	values := []string{tag.Value}
	_, err := StringReflect(field, fieldValue, values)
	if err != nil {
		return err
	}
	return nil
}

func trySetFieldConsumingArgs(field reflect.StructField, fieldValue reflect.Value, tags map[TagType]Tag, args *ArgParserExt) error {
	if !fieldValue.CanSet() {
		return ErrFieldCanNotBeSet
	}

	shortErr := trySetForType(TagShort, field, fieldValue, tags, args)
	longErr := trySetForType(TagLong, field, fieldValue, tags, args)

	if shortErr != nil && longErr != nil {
		return trySetDefault(field, fieldValue, tags)
	}

	return nil
}

// Parse tries to evaluate the given `rawArgs` towards the provided struct `target` (which must include `clapper`-Tags).
// If no `rawArgs` were provided, it defaults to `os.Args[1:]` (all command line arguments without the programm name).
func Parse[T any](target *T, rawArgs ...string) (trailing []string, err error) {
	t := reflect.TypeOf(*target)
	if t.Kind() != reflect.Struct {
		return nil, ErrNoStruct
	}

	if len(rawArgs) == 0 {
		rawArgs = os.Args[1:] // skip the first argument (program name)
	}

	args := NewArgParserExt(rawArgs)

	parsedTags, err := parseStructTags(t)
	if err != nil {
		return nil, err
	}

	reflectValue := reflect.ValueOf(target).Elem()
	processor := NewStructFieldProcessor(t, reflectValue, parsedTags, args)
	for !processor.EOF() {
		if err = processor.Next(); err != nil {
			return nil, err
		}
	}

	if err = processor.Finalize(); err != nil {
		return nil, err
	}

	return processor.GetTrailing(), nil
}
