package gocart

import (
	"net/http"
	"reflect"
	"strings"
)

type EncoderFactory func(writer http.ResponseWriter) Encoder

// Encoder sets data on the return headers and cookies
type Encoder interface {
	Encode(value reflect.Value, field reflect.StructField) error
}

type HeaderEncoder struct {
	headers http.Header
}

func NewHeaderEncoder(writer http.ResponseWriter) Encoder {
	return &HeaderEncoder{headers: writer.Header()}
}

func (enc *HeaderEncoder) Encode(value reflect.Value, field reflect.StructField) error {
	strs, err := EncodePrimitives(value)
	if err != nil {
		return err
	}

	enc.headers.Set(field.Name, strings.Join(strs, ","))
	return nil
}
