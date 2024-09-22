package clapper

import (
	"reflect"
	"strconv"
)

func ptr[T any](t T) *T {
	return &t
}

func ValueFromString(fieldType reflect.Type, inputs []string) (*reflect.Value, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return ptr(reflect.ValueOf(inputs[0])), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.Atoi(inputs[0])
		if err != nil {
			return nil, err
		}
		return ptr(reflect.ValueOf(num)), nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(inputs[0], 64)
		if err != nil {
			return nil, err
		}
		return ptr(reflect.ValueOf(val)), nil
	case reflect.Bool:
		b := true
		var err error
		if len(inputs) > 0 {
			b, err = strconv.ParseBool(inputs[0])
			if err != nil {
				return nil, err
			}
		}
		return ptr(reflect.ValueOf(b)), nil
	default:
		return nil, NewUnsupportedReflectTypeError(fieldType.String())
	}
}

func StringReflect(field reflect.StructField, fieldValue reflect.Value, values []string) ([]string, error) {
	took := make([]string, 0)
	switch field.Type.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(field.Type, len(values), len(values))
		took = values
		for i, value := range values {
			elem := reflect.New(field.Type.Elem()).Elem()
			refValue, err := ValueFromString(field.Type.Elem(), []string{value})
			if err != nil {
				return nil, err
			}
			elem.Set(*refValue)
			slice.Index(i).Set(elem)
		}
		fieldValue.Set(slice)
	case reflect.Pointer:
		elem := reflect.New(field.Type.Elem()).Elem()
		ind := reflect.Indirect(elem)
		v, err := ValueFromString(ind.Type(), values)
		if err != nil {
			return nil, err
		}
		elem.Set(*v)
		fieldValue.Set(elem.Addr())
	default:
		if len(values) > 0 {
			took = append(took, values[0])
		}
		value, err := ValueFromString(field.Type, values)
		if err != nil {
			return nil, err
		}
		fieldValue.Set(*value)
	}

	return took, nil
}
