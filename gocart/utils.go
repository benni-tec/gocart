package gocart

import (
	"net/http"
	"reflect"
)

func genericToType[T any]() reflect.Type {
	var array []T
	return reflect.TypeOf(array).Elem()
}

func setHeader[T any](header http.Header, key string, values Serializer[T]) {
	if values == nil {
		return
	}

	header.Del(key)

	for _, value := range values.Type().HttpType {
		header.Add(key, value)
	}
}
