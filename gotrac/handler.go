package gotrac

import (
	"net/http"
	"reflect"
)

// Handler is just a http.Handler which can also provide HandlerInformation
type Handler interface {
	http.Handler
	Info() *HandlerInformation
}

// HandlerType defines both the type used in go and the http content type (MIME-type)
type HandlerType struct {
	GoType   reflect.Type
	HttpType []string
}
