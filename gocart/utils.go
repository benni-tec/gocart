package gocart

import "reflect"

func genericToType[T any]() reflect.Type {
	var array []T
	return reflect.TypeOf(array).Elem()
}
