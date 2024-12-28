package gocart

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type DecoderFactory func(request *http.Request) Decoder

// Decoder reads data from HTTP headers, the query, formdata, etc.
// The returned strings are converted to the proper primitive types.
// If the receiving field is not an array, but multiple values are provided the first one will be used
type Decoder interface {
	Decode(field reflect.StructField) ([]string, error)
}

type HeaderDecoder struct {
	headers http.Header
}

func NewHeaderDecoder(request *http.Request) Decoder {
	return &HeaderDecoder{headers: request.Header}
}

func (dec *HeaderDecoder) Decode(field reflect.StructField) ([]string, error) {
	if tag, ok := field.Tag.Lookup("meta"); ok {
		return strings.Split(dec.headers.Get(tag), ","), nil
	}

	return nil, nil
}

type UrlValuesDecoder struct {
	name   string
	values url.Values
}

func NewQueryDecoder(request *http.Request) Decoder {
	return &UrlValuesDecoder{name: "query", values: request.URL.Query()}
}

func NewFormDecorder(request *http.Request) Decoder {
	return &UrlValuesDecoder{name: "form", values: request.Form}
}

func (dec *UrlValuesDecoder) Decode(field reflect.StructField) ([]string, error) {
	if tag, ok := field.Tag.Lookup(dec.name); ok {
		return strings.Split(dec.values.Get(tag), ","), nil
	}

	return nil, nil
}

type PathDecoder struct {
	request *http.Request
}

func NewPathDecoder(request *http.Request) Decoder {
	return &PathDecoder{request: request}
}

func (dec *PathDecoder) Decode(field reflect.StructField) ([]string, error) {
	if tag, ok := field.Tag.Lookup("path"); ok {
		return []string{dec.request.PathValue(tag)}, nil
	}

	return nil, nil
}
