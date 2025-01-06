package gocrew

import (
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
)

// Generator generates a T documentation from a chi.Router
type Generator[T any] interface {
	Generate(router chi.Routes) (*T, error)
}

// OpenApi31 returns a Generator that generates a OpenAPI v3.1.0 compliant specification
func OpenApi31(title string, defaultTag *openapi31.Tag) Generator[OpenApi31Spec] {
	if defaultTag == nil {
		defaultTag = &openapi31.Tag{
			Name:        "default",
			Description: P("This controller catches all endpoints that are not beneath a controller!"),
		}
	}

	return &openapi31Generator{
		title:            title,
		defaultTag:       *defaultTag,
		defaultResponses: map[string]openapi31.ResponseOrReference{},
	}
}
