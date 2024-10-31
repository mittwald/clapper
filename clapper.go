package clapper

import (
	"os"
	"reflect"
	"strings"
)

const (
	TagName = "clapper"
)

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

func parseStructTags(t reflect.Type) (map[int]map[TagType]Tag, error) {
	parsedTags := make(map[int]map[TagType]Tag, 0)
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
		parsedTags[i] = tags
	}
	return parsedTags, nil
}

func getValues(tag Tag, args *ArgsParser) []string {
	lookup := tag.Name
	if tag.HasValue() {
		lookup = tag.Value
	}

	values, ok := args.Params[lookup]

	if !ok {
		return nil
	}

	return values
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

func valuesFor(tagType TagType, tags map[TagType]Tag, args *ArgsParser) []string {
	tag, ok := tags[tagType]
	if !ok {
		return nil
	}
	values := getValues(tag, args)
	return values
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

func trySetField(field reflect.StructField, fieldValue reflect.Value, tags map[TagType]Tag, args *ArgsParser) error {
	if !fieldValue.CanSet() {
		return ErrFieldCanNotBeSet
	}

	defaulted := false
	values := valuesFor(TagLong, tags, args)
	if values != nil {
		// Ugly way to get rid of exfra short value
		temp := valuesFor(TagShort, tags, args)
		if temp != nil {
			args.PopTrailing(1)
		}
	} else {
		values = valuesFor(TagShort, tags, args)
		if values == nil {
			tag, ok := tags[TagDefault]
			if !ok {
				// Pointers and bools are optional by default
				if isOptionalField(field) {
					return nil
				}
				return NewMandatoryParameterError(paramName(tags))
			}
			values = []string{tag.Value}
			defaulted = true
		}
	}

	took, err := StringReflect(field, fieldValue, values)
	if err != nil {
		return err
	}

	if defaulted {
		return nil
	}

	args.PopTrailing(took)

	return nil
}

func Parse[T any](target *T, rawArgs ...string) (trailing []string, err error) {
	t := reflect.TypeOf(*target)
	if t.Kind() != reflect.Struct {
		return nil, ErrNoStruct
	}

	if len(rawArgs) == 0 {
		rawArgs = os.Args[1:] // skip the first argument (program name)
	}

	args, err := NewArgsParser().Parse(rawArgs)
	if err != nil {
		return nil, err
	}

	parsedTags, err := parseStructTags(t)
	if err != nil {
		return nil, err
	}

	reflectValue := reflect.ValueOf(target).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := reflectValue.Field(i)

		tags, ok := parsedTags[i]
		if !ok {
			continue
		}

		err = trySetField(field, fieldValue, tags, args)
		if err != nil {
			return nil, err
		}
	}

	return args.Trailing, nil
}
