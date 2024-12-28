package gocart

import (
	"encoding/json"
	"encoding/xml"
	"github.com/benni-tec/gocart/gotrac"
	"gopkg.in/yaml.v3"
	"net/http"
	"reflect"
)

// Serializer determines how the body of the http.Request should be decoded.
// Common types can be constructed using Json, Yaml, Xml, Binary, etc.
// However, the Serializer can also be implemented by you!
type Serializer[T any] interface {
	Serialize(body *T, headers http.Header) ([]byte, error)
	Deserialize(data []byte, headers http.Header) (*T, error)
	Type() *gotrac.HandlerType
}

// +++ JSON, YAML +++

// MarshalSerializer implements the Serializer interface using the go-convention marshal/unmarshal functions.
// It is for example used for Json, Yaml and Xml
type MarshalSerializer[T any] struct {
	marshal   func(v any) ([]byte, error)
	unmarshal func(data []byte, v any) error
	mimeTypes []string
}

// Json Serializer to decode the http.Request`s body
func Json[T any]() Serializer[T] {
	return &MarshalSerializer[T]{
		marshal:   json.Marshal,
		unmarshal: json.Unmarshal,
		mimeTypes: []string{"application/json"},
	}
}

// Yaml Serializer to decode the http.Request`s body
func Yaml[T any]() Serializer[T] {
	return &MarshalSerializer[T]{
		marshal:   yaml.Marshal,
		unmarshal: yaml.Unmarshal,
		mimeTypes: []string{"application/x-yaml", "text/yaml"},
	}
}

// Xml Serializer to decode the http.Request`s body
func Xml[T any]() Serializer[T] {
	return &MarshalSerializer[T]{
		marshal:   xml.Marshal,
		unmarshal: xml.Unmarshal,
		mimeTypes: []string{"application/xml"},
	}
}

func (j *MarshalSerializer[T]) Serialize(body *T, headers http.Header) ([]byte, error) {
	return json.Marshal(body)
}

func (j *MarshalSerializer[T]) Deserialize(data []byte, headers http.Header) (*T, error) {
	value := new(T)
	err := json.Unmarshal(data, value)
	return value, err
}

func (j *MarshalSerializer[T]) Type() *gotrac.HandlerType {
	return &gotrac.HandlerType{
		GoType:   genericToType[T](),
		HttpType: []string{"application/json"},
	}
}

// +++ Binary +++

type BinarySerializer struct {
	contentType []string
}

// Binary Serializer just passes through the body, but it does allow you to set the MIME type yourself!
func Binary(contentType ...string) Serializer[[]byte] {
	return &BinarySerializer{contentType: contentType}
}

func (receiver *BinarySerializer) Serialize(body *[]byte, headers http.Header) ([]byte, error) {
	return *body, nil
}

func (receiver *BinarySerializer) Deserialize(data []byte, headers http.Header) (*[]byte, error) {
	return &data, nil
}

func (receiver *BinarySerializer) Type() *gotrac.HandlerType {
	return &gotrac.HandlerType{
		GoType:   reflect.TypeOf([]byte{}),
		HttpType: receiver.contentType,
	}
}
