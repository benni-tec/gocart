package gotrac

import (
	"github.com/benni-tec/gocart/goflag"
	middleware2 "github.com/benni-tec/gocart/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// Mux implements the Router using a chi.Router
type Mux struct {
	router chi.Router
	info   goflag.Information
}

// Default creates a new Router and adds common middlewares, that should be present on the root router
func Default() Router {
	mux := NewRouter()
	mux.Use(middleware.RedirectSlashes)
	mux.Use(middleware.Logger)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware2.ErrorMiddleware)

	return mux
}

// NewRouter creates a new Router without any middlewares.
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

func (m *Mux) Handle(pattern string, h http.Handler) Route {
	actor := wrapToHandler(h)
	m.router.Handle(pattern, actor)
	return actor
}

func (m *Mux) HandleFunc(pattern string, h http.HandlerFunc) Route {
	actor := wrapFuncToHandler(h)
	m.router.Handle(pattern, actor)
	return actor
}

func (m *Mux) Method(method, pattern string, h http.Handler) Route {
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

func (m *Mux) Info() *goflag.Information {
	return &m.info
}

func (m *Mux) WithInfo(fn func(info *goflag.Information)) Router {
	if fn != nil {
		fn(&m.info)
	}

	return m
}
