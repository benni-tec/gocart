package gocart

import (
	"gocart/gotrac"
	"net/http"
)

type Empty struct{}

func None() Converter[Empty] {
	return &EmptyConverter{}
}

type EmptyConverter struct{}

func (e EmptyConverter) Serialize(body *Empty, headers http.Header) ([]byte, error) {
	return []byte{}, nil
}

func (e EmptyConverter) Deserialize(data []byte, headers http.Header) (*Empty, error) {
	return &Empty{}, nil
}

func (e EmptyConverter) Type() *gotrac.HandlerType {
	return nil
}
