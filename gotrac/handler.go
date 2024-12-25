package gotrac

import (
	"net/http"
	"reflect"
)

type Handler interface {
	http.Handler
	Info() *HandlerInformation
}

type HandlerType struct {
	GoType   reflect.Type
	HttpType []string
}
