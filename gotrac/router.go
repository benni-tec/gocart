package gotrac

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Router consisting of the gotrac routing methods used by chi's Mux,
// using only the standard net/http.
type Router interface {
	http.Handler
	chi.Routes

	// Use appends one or more middlewares onto the Router stack.
	Use(middlewares ...func(http.Handler) http.Handler)

	// With adds inline middlewares for an endpoint handler.
	With(middlewares ...func(http.Handler) http.Handler) Router

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-Router.
	Group(fn func(r Router)) Router

	// Route mounts a sub-Router along a `pattern`` string.
	Route(pattern string, fn func(r Router)) Router

	// Mount attaches another http.Handler along ./pattern/*
	Mount(pattern string, h http.Handler)

	// Handle and HandleFunc adds routes for `pattern` that matches
	// all HTTP methods.
	Handle(pattern string, h Handler) Route
	HandleFunc(pattern string, h http.HandlerFunc) Route

	// Method and MethodFunc adds routes for `pattern` that matches
	// the `method` HTTP method.
	Method(method, pattern string, h Handler) Route
	MethodFunc(method, pattern string, h http.HandlerFunc) Route

	// NotFound defines a handler to respond whenever a route could
	// not be found.
	NotFound(h http.HandlerFunc)

	// MethodNotAllowed defines a handler to respond whenever a method is
	// not allowed.
	MethodNotAllowed(h http.HandlerFunc)

	// ++++ Information ++++
	Info() *RouterInformation
	WithInfo(fn func(info *RouterInformation)) Router
}
