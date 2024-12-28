package gotrac

import "net/http"

// Route is a handler that has been registered to a Router.
// When registering with a Router any handler is wrapped to a Route therefore allowing the information to be edited.
type Route interface {
	Handler
	WithInfo(fn func(route *RouteInformation)) Route
}

type routeImpl struct {
	handler http.HandlerFunc
	info    RouteInformation
}

func wrapToHandler(handler Handler) *routeImpl {
	return wrapToActorFunc(handler.Info(), func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
}

func wrapToActorFunc(info *HandlerInformation, handler http.HandlerFunc) *routeImpl {
	return &routeImpl{
		handler: handler,
		info:    RouteInformation(*info),
	}
}

func wrapFuncToHandler(handler http.HandlerFunc) *routeImpl {
	return &routeImpl{
		handler: handler,
		info: RouteInformation{
			Information: Information{
				Summary:     "",
				Description: "",
			},
			Input:  nil,
			Output: nil,
			Hidden: false,
		},
	}
}

func (a *routeImpl) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	a.handler(writer, request)
}

func (a *routeImpl) Info() *HandlerInformation {
	cast := HandlerInformation(a.info)
	return &cast
}

func (a *routeImpl) WithInfo(fn func(route *RouteInformation)) Route {
	if fn != nil {
		fn(&a.info)
	}

	return a
}
