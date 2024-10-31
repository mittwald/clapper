package clapper

import (
	"reflect"
	"strconv"
)

func ptr[T any](t T) *T {
	return &t
}

func ValueFromString(fieldType reflect.Type, inputs []string) (*reflect.Value, int, error) {
	switch fieldType.Kind() {
	case reflect.String:
		return ptr(reflect.ValueOf(inputs[0])), 1, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num, err := strconv.Atoi(inputs[0])
		if err != nil {
			return nil, 0, err
		}
		return ptr(reflect.ValueOf(num)), 1, nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(inputs[0], 64)
		if err != nil {
			return nil, 0, err
		}
		return ptr(reflect.ValueOf(val)), 1, nil
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
			refValue, tookCount, err := ValueFromString(field.Type.Elem(), []string{value})
			if err != nil {
				return 0, err
			}
			took += tookCount
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
