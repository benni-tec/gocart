package gotrac

import "reflect"

// Json returns a HandlerType that can be used to represent a json request/response body that will be serialized to T.
func Json[T any]() *HandlerType {
	return &HandlerType{
		GoType:   genericToType[T](),
		HttpType: []string{"application/json"},
	}
}

// Yaml returns a HandlerType that can be used to represent a yaml request/response body that will be serialized to T.
func Yaml[T any]() *HandlerType {
	return &HandlerType{
		GoType:   genericToType[T](),
		HttpType: []string{"application/x-yaml", "text/yaml"},
	}
}

// Xml returns a HandlerType that can be used to represent a xml request/response body that will be serialized to T.
func Xml[T any]() *HandlerType {
	return &HandlerType{
		GoType:   genericToType[T](),
		HttpType: []string{"application/xml"},
	}
}

// File returns a HandlerType that can be used to represent a request/response body that is a file.
func File() *HandlerType {
	return Binary("application/octet-stream")
}

// Binary returns a HandlerType that represents a binary request/response body.
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
