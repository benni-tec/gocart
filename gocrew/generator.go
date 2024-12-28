package gocrew

import (
	"github.com/benni-tec/gocart/gotrac"
	"github.com/swaggest/openapi-go/openapi31"
)

// Generator generates a T documentation from a gotrac.Router
type Generator[T any] interface {
	Generate(router gotrac.Router) (*T, error)
}

// OpenApi31 returns a Generator that generates a OpenAPI v3.1.0 compliant specification
func OpenApi31() Generator[OpenApi31Spec] {
	return &openapi31Generator{
		reflector:        openapi31.NewReflector(),
		defaultResponses: map[string]openapi31.ResponseOrReference{},
	}
}
