package gotrac

import "reflect"

func Json[T any]() *HandlerType {
	return &HandlerType{
		GoType:   genericToType[T](),
		HttpType: []string{"application/json"},
	}
}

func File() *HandlerType {
	return Binary("application/octet-stream")
}

func Binary(mimeTypes ...string) *HandlerType {
	return &HandlerType{
		GoType:   reflect.TypeOf([]byte{}),
		HttpType: mimeTypes,
	}
}

func genericToType[T any]() reflect.Type {
	var array []T
	return reflect.TypeOf(array).Elem()
}
