package gotrac

import (
	"github.com/benni-tec/gocart/middleware"
	"github.com/go-chi/chi/v5"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"net/http"
)

type Mux struct {
	router chi.Router

	summary     string
	description string
}

// Default creates a new Router and adds the IdMiddleware required for error handling
func Default() Router {
	mux := wrapToRouter(chi.NewRouter())
	mux.Use(middleware.IdMiddleware)
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
	return wrapToRouter(m.router.Group(func(r chi.Router) {
		fn(wrapToRouter(m.router))
	}))
}

func (m *Mux) Route(pattern string, fn func(r Router)) Router {
	return wrapToRouter(m.router.Route(pattern, func(r chi.Router) {
		fn(wrapToRouter(m.router))
	}))
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

// +++ Docs +++

func (m *Mux) WithDocs(pattern string, generator Generator) {
	docs, err := generator.Generate(m)
	if err != nil {
		panic(err)
	}

	m.router.Method(http.MethodGet, pattern, specHandler(docs))
}

func (m *Mux) WithSwaggerUI(pattern string, docPattern string, title string, config *swg.Config) {
	var handler http.Handler
	if config != nil {
		ui := swgui.NewWithConfig(*config)
		handler = ui(title, docPattern, pattern)
	} else {
		handler = swgui.New(title, docPattern, pattern)
	}

	m.Mount(pattern, handler)
}

// +++ Information +++

func (m *Mux) Summary() string {
	return m.summary
}

func (m *Mux) Description() string {
	return m.description
}

func (m *Mux) WithSummary(title string) Router {
	m.summary = title
	return m
}

func (m *Mux) WithDescription(description string) Router {
	m.description = description
	return m
}
