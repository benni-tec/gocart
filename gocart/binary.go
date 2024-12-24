package gocart

import (
	"gocart/gotrac"
	"net/http"
	"reflect"
)

type BinaryConverter struct {
	contentType []string
}

func Binary(contentType ...string) Converter[[]byte] {
	return &BinaryConverter{contentType: contentType}
}

func (receiver *BinaryConverter) Serialize(body *[]byte, headers http.Header) ([]byte, error) {
	return *body, nil
}

func (receiver *BinaryConverter) Deserialize(data []byte, headers http.Header) (*[]byte, error) {
	return &data, nil
}

func (receiver *BinaryConverter) Type() *gotrac.HandlerType {
	return &gotrac.HandlerType{
		GoType:   reflect.TypeOf([]byte{}),
		HttpType: receiver.contentType,
	}
}
