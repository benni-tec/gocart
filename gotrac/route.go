package gotrac

import (
	"github.com/benni-tec/gocart/goflag"
	"net/http"
)

// Route is a handler that has been registered to a Router.
// When registering with a Router any handler is wrapped to a Route therefore allowing the information to be edited.
type Route interface {
	goflag.EndpointFlag
	WithInfo(fn func(route *goflag.EndpointInformation)) Route
}

type routeImpl struct {
	handler http.HandlerFunc
	info    goflag.EndpointInformation
}

func wrapToHandler(handler http.Handler) *routeImpl {
	var info *goflag.EndpointInformation = nil
	if with, ok := handler.(goflag.EndpointFlag); ok {
		info = with.Info()
	}

	return wrapToActorFunc(info, func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
}

func wrapToActorFunc(info *goflag.EndpointInformation, handler http.HandlerFunc) *routeImpl {
	if info == nil {
		info = &goflag.EndpointInformation{}
	}

	return &routeImpl{
		handler: handler,
		info:    *info,
	}
}

func wrapFuncToHandler(handler http.HandlerFunc) *routeImpl {
	return &routeImpl{
		handler: handler,
		info: goflag.EndpointInformation{
			Information: goflag.Information{
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

func (a *routeImpl) Info() *goflag.EndpointInformation {
	cast := goflag.EndpointInformation(a.info)
	return &cast
}

func (a *routeImpl) WithInfo(fn func(route *goflag.EndpointInformation)) Route {
	if fn != nil {
		fn(&a.info)
	}

	return a
}
