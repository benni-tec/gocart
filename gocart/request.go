package gocart

import (
	"net/http"
)

// Request is a http.Request that also contains the already deserialized body,
// which can then be retrieved by calling Body().
type Request[TBody any] struct {
	http.Request
	body *TBody
}

func wrapToBodyRequest[TBody any](request *http.Request, body *TBody) *Request[TBody] {
	return &Request[TBody]{
		Request: *request,
		body:    body,
	}
}

// Body returns a pointer to the already deserialized body.
func (r *Request[TBody]) Body() *TBody {
	return r.body
}
