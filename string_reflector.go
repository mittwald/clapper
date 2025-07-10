package clapper

import (
	"fmt"
	"reflect"
	"strconv"

	internalerrors "github.com/mittwald/clapper/internal/errors"
)

func ptr[T any](t T) *T {
	return &t
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

func trySetForType(
	tagType TagType,
	field reflect.StructField,
	fieldValue reflect.Value,
	tags map[TagType]Tag,
	args *ArgParserExt,
) error {
	key, values := valuesFor(tagType, tags, args)
	if values == nil {
		return internalerrors.ErrInternalNoArgumentsForTag
	}

	took, err := StringReflect(field, fieldValue, values)
	if err != nil {
		return err
	}

	argType := mustTagTypeToArgType(tagType)
	args.Consume(key, argType, took)

	return nil
}

func trySetDefault(field reflect.StructField, fieldValue reflect.Value, tags TagMap) error {
	tag, ok := tags[TagDefault]
	if !ok {
		if isOptionalField(field) {
			return nil
		}
		return NewMandatoryParameterError(tags.InputArgument())
	}
	values := []string{tag.Value}
	_, err := StringReflect(field, fieldValue, values)
	if err != nil {
		return err
	}
	return nil
}

func trySetFieldConsumingArgs(
	field reflect.StructField,
	fieldValue reflect.Value,
	tags TagMap,
	args *ArgParserExt,
) error {
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

func inputNeededForKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool:
		return false
	default:
		return true
	}
}

func parseFloat(input string, bits int) (*reflect.Value, error) {
	val, err := strconv.ParseFloat(input, bits)
	if err != nil {
		return nil, err
	}

	switch bits {
	case 32:
		val32 := float32(val)
		return ptr(reflect.ValueOf(val32)), nil
	case 64:
		return ptr(reflect.ValueOf(val)), nil
	}

	return nil, NewUnsupportedReflectTypeError(fmt.Sprintf("float%d", bits))
}

func ValueFromString(fieldType reflect.Type, inputs []string) (*reflect.Value, int, error) {
	if inputNeededForKind(fieldType.Kind()) && len(inputs) == 0 {
		return nil, 0, ErrEmptyArgument
	}

	switch fieldType.Kind() {
	case reflect.String:
		return ptr(reflect.ValueOf(inputs[0])), 1, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.Atoi(inputs[0])
		if err != nil {
			err = NewUnexpectedInputFormatError(inputs[0], fieldType)
			return nil, 0, err
		}
		return ptr(reflect.ValueOf(num)), 1, nil
	case reflect.Float32, reflect.Float64:
		val, err := parseFloat(inputs[0], fieldType.Bits())
		if err != nil {
			err = NewUnexpectedInputFormatError(inputs[0], fieldType)
			return nil, 0, err
		}
		return val, 1, nil
	case reflect.Bool:
		b := true
		return ptr(reflect.ValueOf(b)), 0, nil
	default:
		return nil, 0, NewUnsupportedReflectTypeError(fieldType.String())
	}
}

func StringReflect(field reflect.StructField, fieldValue reflect.Value, values []string) (int, error) {
	took := 0
	switch field.Type.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(field.Type, len(values), len(values))
		took = len(values)
		for i, value := range values {
			elem := reflect.New(field.Type.Elem()).Elem()
			refValue, _, err := ValueFromString(field.Type.Elem(), []string{value})
			if err != nil {
				return 0, err
			}
			elem.Set(*refValue)
			slice.Index(i).Set(elem)
		}
		fieldValue.Set(slice)
	case reflect.Pointer:
		elem := reflect.New(field.Type.Elem()).Elem()
		ind := reflect.Indirect(elem)
		v, tookCount, err := ValueFromString(ind.Type(), values)
		if err != nil {
			return 0, err
		}
		took += tookCount
		elem.Set(*v)
		fieldValue.Set(elem.Addr())
	default:
		value, tookCount, err := ValueFromString(field.Type, values)
		if err != nil {
			return 0, err
		}
		took += tookCount
		fieldValue.Set(*value)
	}

	return took, nil
}
