package gocrew

import (
	"github.com/benni-tec/gocart/goflag"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

type WalkFunc func(method string, route string, handler http.Handler, controller goflag.ControllerFlag) error
type ControllerFunc func(controller goflag.ControllerFlag) error

func Walk(r chi.Routes, walkFn WalkFunc, onController ControllerFunc) error {
	return walk(r, walkFn, onController, "", nil)
}

// copied from chi, removed middlewares, added controllers
func walk(r chi.Routes, walkFn WalkFunc, onController ControllerFunc, parentRoute string, controller goflag.ControllerFlag) error {
	for _, route := range r.Routes() {
		if route.SubRoutes != nil {
			current := controller
			if cont, ok := route.SubRoutes.(goflag.ControllerFlag); ok {
				err := onController(cont)
				if err != nil {
					return err
				}

				current = cont
			}

			if err := walk(route.SubRoutes, walkFn, onController, parentRoute+route.Pattern, current); err != nil {
				return err
			}

			continue
		}

		for method, handler := range route.Handlers {
			if method == "*" {
				// Ignore a "catchAll" method, since we pass down all the specific methods for each route.
				continue
			}

			fullRoute := parentRoute + route.Pattern
			fullRoute = strings.Replace(fullRoute, "/*/", "/", -1)

			if chain, ok := handler.(*chi.ChainHandler); ok {
				if err := walkFn(method, fullRoute, chain.Endpoint, controller); err != nil {
					return err
				}
			} else {
				if err := walkFn(method, fullRoute, handler, controller); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
