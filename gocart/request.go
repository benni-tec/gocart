package gocart

import (
	"net/http"
)

type Request[TBody any] struct {
	http.Request
	body *TBody
}

func ConvertRequest[TBody any](request *http.Request, body *TBody) *Request[TBody] {
	return &Request[TBody]{
		Request: *request,
		body:    body,
	}
}

func (r *Request[TBody]) Body() *TBody {
	return r.body
}
