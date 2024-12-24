package gocart

import (
	"encoding/json"
	"github.com/benni-tec/gocart/gotrac"
	"net/http"
)

type JsonConverter[T any] struct {
}

func Json[T any]() Converter[T] {
	return &JsonConverter[T]{}
}

func (j JsonConverter[T]) Serialize(body *T, headers http.Header) ([]byte, error) {
	return json.Marshal(body)
}

func (j JsonConverter[T]) Deserialize(data []byte, headers http.Header) (*T, error) {
	value := new(T)
	err := json.Unmarshal(data, value)
	return value, err
}

func (j JsonConverter[T]) Type() *gotrac.HandlerType {
	return &gotrac.HandlerType{
		GoType:   genericToType[T](),
		HttpType: []string{"application/json"},
	}
}
