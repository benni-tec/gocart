package gocart

import (
	"net/http"
	"reflect"
)

func genericToType[T any]() reflect.Type {
	return reflect.TypeOf(new(T))
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
