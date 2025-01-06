package gotrac

import (
	"github.com/benni-tec/gocart/goflag"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func WithName(name string, handler goflag.Controller) goflag.ControllerFlag {
	return &controllerImpl{
		handler: handler,
		info: goflag.ControllerInformation{
			Name: name,
		},
	}
}

type controllerImpl struct {
	handler goflag.Controller
	info    goflag.ControllerInformation
}

func (c *controllerImpl) Routes() []chi.Route {
	return c.handler.Routes()
}

func (c *controllerImpl) Middlewares() chi.Middlewares {
	return c.handler.Middlewares()
}

func (c *controllerImpl) Match(rctx *chi.Context, method, path string) bool {
	return c.handler.Match(rctx, method, path)
}

func (c *controllerImpl) Find(rctx *chi.Context, method, path string) string {
	return c.handler.Find(rctx, method, path)
}

func (c *controllerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.handler.ServeHTTP(w, r)
}

func (c *controllerImpl) Info() *goflag.ControllerInformation {
	return &c.info
}
