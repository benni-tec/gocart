package gotrac

import (
	"net/http"
	"reflect"
)

type Handler interface {
	http.Handler
	Info() HandlerInformation
}

type HandlerType struct {
	GoType   reflect.Type
	HttpType []string
}

type handlerImpl struct {
	handler http.HandlerFunc

	summary     string
	description string
	input       *HandlerType
	output      *HandlerType
	hidden      bool
}

func wrapToHandler(handler Handler) *handlerImpl {
	return wrapToActorFunc(handler.Info(), func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
}

func wrapToActorFunc(info HandlerInformation, handler http.HandlerFunc) *handlerImpl {
	return &handlerImpl{
		handler:     handler,
		summary:     info.Summary(),
		description: info.Description(),
		input:       info.Input(),
		output:      info.Output(),
		hidden:      false,
	}
}

func wrapFuncToHandler(handler http.HandlerFunc) *handlerImpl {
	return &handlerImpl{
		handler:     handler,
		summary:     "",
		description: "",
		input:       nil,
		output:      nil,
		hidden:      false,
	}
}

func _ensureRoute() Route {
	return &handlerImpl{}
}

func _ensureHandler() Handler {
	return &handlerImpl{}
}

// Handler

func (a *handlerImpl) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.handler(writer, request)
}

func (a *handlerImpl) Info() HandlerInformation {
	return a
}

// HandlerInformation: get data

func (a *handlerImpl) Summary() string {
	return a.summary
}

func (a *handlerImpl) Description() string {
	return a.description
}

func (a *handlerImpl) Input() *HandlerType {
	return a.input
}

func (a *handlerImpl) Output() *HandlerType {
	return a.output
}

func (a *handlerImpl) Hidden() bool {
	return a.hidden
}

// Route: Fluent API to set data

func (a *handlerImpl) WithSummary(title string) Route {
	a.summary = title
	return a
}

func (a *handlerImpl) WithDescription(description string) Route {
	a.description = description
	return a
}

func (a *handlerImpl) ForInput(fluent func(typ *HandlerType)) Route {
	fluent(a.input)
	return a
}

func (a *handlerImpl) ForOutput(fluent func(typ *HandlerType)) Route {
	fluent(a.output)
	return a
}

func (a *handlerImpl) WithInput(typ *HandlerType) Route {
	a.input = typ
	return a
}

func (a *handlerImpl) WithOutput(typ *HandlerType) Route {
	a.output = typ
	return a
}

func (a *handlerImpl) WithHidden(hidden bool) Route {
	a.hidden = hidden
	return a
}
