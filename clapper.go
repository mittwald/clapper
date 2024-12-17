package clapper

import (
	"fmt"
	"os"
	"reflect"
)

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

func valuesFor(tagType TagType, tags TagMap, args *ArgParserExt) (key string, values []string) {
	tag, ok := tags[tagType]
	if !ok {
		return "", nil
	}
	argType := mustTagTypeToArgType(tag.Type)

	key = tag.ArgumentName()

	values, ok = args.Get(key, argType)
	if !ok {
		return key, nil
	}
	return key, values
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
