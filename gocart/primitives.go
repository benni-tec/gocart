package gocart

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// EncodePrimitives encodes the value(s) of value as a string and if multiple separates them by a "," (see url.Values)
func EncodePrimitives(value reflect.Value) ([]string, error) {
	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		str, err := EncodePrimitive(value)
		if err != nil {
			return nil, err
		}

		return []string{str}, nil
	}

	var strings []string
	for i := range value.Len() {
		str, err := EncodePrimitive(value.Index(i))
		if err != nil {
			return nil, err
		}

		strings = append(strings, str)
	}

	return strings, nil

}

// EncodePrimitive encodes the value of value as a string
func EncodePrimitive(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, 64), nil
	case reflect.String:
		return value.String(), nil
	default:
		return "", errors.New(fmt.Sprintf("gocart: %s is not a primitive type", value.Kind()))
	}
}

// AssignPrimitives is used when decoding, to convert the strs to the appropriate primitives and then assign them to value
func AssignPrimitives(value reflect.Value, strs []string) error {
	if len(strs) == 0 {
		return nil
	}

	if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
		return AssignPrimitive(value, strs[0])
	}

	value.SetLen(len(strs))
	for i, str := range strs {
		err := AssignPrimitive(value.Index(i), str)
		if err != nil {
			return err
		}
	}

	return nil
}

// AssignPrimitive is used when decoding, to convert the str to the appropriate primitive and then assign it to value
func AssignPrimitive(value reflect.Value, str string) error {
	switch value.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}

		value.SetBool(b)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}

		value.SetInt(i)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
		}

		value.SetUint(u)
		return nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}

		value.SetFloat(f)
		return nil
	case reflect.String:
		value.SetString(str)
		return nil
	default:
		return errors.New(fmt.Sprintf("gocart: %s is not a primitive type", value.Kind()))
	}
}
