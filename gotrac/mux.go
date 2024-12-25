package gotrac

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Mux struct {
	router chi.Router
	info   RouterInformation
}

// Default creates a new Router and adds the IdMiddleware required for error handling
func Default() Router {
	mux := NewRouter()
	mux.Use(middleware.RequestID)
	return mux
}

func NewRouter() Router {
	mux := wrapToRouter(chi.NewRouter())
	return mux
}

func wrapToRouter(r chi.Router) Router {
	return &Mux{
		router: r,
	}
}

func (m *Mux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.router.ServeHTTP(writer, request)
}

func (m *Mux) Routes() []chi.Route {
	return m.router.Routes()
}

func (m *Mux) Middlewares() chi.Middlewares {
	return m.router.Middlewares()
}

func (m *Mux) Match(rctx *chi.Context, method, path string) bool {
	return m.router.Match(rctx, method, path)
}

func (m *Mux) Find(rctx *chi.Context, method, path string) string {
	return m.router.Find(rctx, method, path)
}

func (m *Mux) Use(middlewares ...func(http.Handler) http.Handler) {
	m.router.Use(middlewares...)
}

func (m *Mux) With(middlewares ...func(http.Handler) http.Handler) Router {
	return wrapToRouter(m.router.With(middlewares...))
}

func (m *Mux) Group(fn func(r Router)) Router {
	inline := m.With()

	if fn != nil {
		fn(inline)
	}

	return inline
}

func (m *Mux) Route(pattern string, fn func(r Router)) Router {
	sub := NewRouter()

	if fn != nil {
		fn(sub)
	}

	m.Mount(pattern, sub)
	return sub
}

func (m *Mux) Mount(pattern string, h http.Handler) {
	m.router.Mount(pattern, h)
}

func (m *Mux) Handle(pattern string, h Handler) Route {
	actor := wrapToHandler(h)
	m.router.Handle(pattern, actor)
	return actor
}

func (m *Mux) HandleFunc(pattern string, h http.HandlerFunc) Route {
	actor := wrapFuncToHandler(h)
	m.router.Handle(pattern, actor)
	return actor
}

func (m *Mux) Method(method, pattern string, h Handler) Route {
	actor := wrapToHandler(h)
	m.router.Method(method, pattern, actor)
	return actor
}

func (m *Mux) MethodFunc(method, pattern string, h http.HandlerFunc) Route {
	actor := wrapFuncToHandler(h)
	m.router.Method(method, pattern, actor)
	return actor
}

func (m *Mux) NotFound(h http.HandlerFunc) {
	m.router.NotFound(h)
}

func (m *Mux) MethodNotAllowed(h http.HandlerFunc) {
	m.router.MethodNotAllowed(h)
}

// +++ Information +++

func (m *Mux) Info() *RouterInformation {
	return &m.info
}

func (m *Mux) WithInfo(fn func(info *RouterInformation)) Router {
	if fn != nil {
		fn(&m.info)
	}

	return m
}
