package gotrac

import (
	"github.com/benni-tec/gocart/goflag"
	"reflect"
)

// Json returns a Type that can be used to represent a json request/response body that will be serialized to T.
func Json[T any]() *goflag.Type {
	return &goflag.Type{
		GoType:   genericToType[T](),
		HttpType: []string{"application/json"},
	}
}

// Yaml returns a Type that can be used to represent a yaml request/response body that will be serialized to T.
func Yaml[T any]() *goflag.Type {
	return &goflag.Type{
		GoType:   genericToType[T](),
		HttpType: []string{"application/x-yaml", "text/yaml"},
	}
}

// Xml returns a Type that can be used to represent a xml request/response body that will be serialized to T.
func Xml[T any]() *goflag.Type {
	return &goflag.Type{
		GoType:   genericToType[T](),
		HttpType: []string{"application/xml"},
	}
}

// File returns a Type that can be used to represent a request/response body that is a file.
func File() *goflag.Type {
	return Binary("application/octet-stream")
}

// Binary returns a Type that represents a binary request/response body.
func Binary(mimeTypes ...string) *goflag.Type {
	return &goflag.Type{
		GoType:   reflect.TypeOf([]byte{}),
		HttpType: mimeTypes,
	}
}

// None is used to declare a type that does not produce a body, but does contain path, query, etc. fields
func None[T any]() *goflag.Type {
	return &goflag.Type{
		GoType:   genericToType[T](),
		HttpType: []string{},
	}
}

func genericToType[T any]() reflect.Type {
	var array []T
	return reflect.TypeOf(array).Elem()
}
