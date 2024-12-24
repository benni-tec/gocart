package gotrac

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
	"path"
	"reflect"
	"strings"
)

type Generator interface {
	Generate(router Router) (*openapi31.Spec, error)
}

type generatorImpl struct {
	reflector        *openapi31.Reflector
	defaultResponses map[string]openapi31.ResponseOrReference
}

func NewGenerator() Generator {
	return &generatorImpl{
		reflector:        openapi31.NewReflector(),
		defaultResponses: map[string]openapi31.ResponseOrReference{},
	}
}

func (g *generatorImpl) Generate(router Router) (*openapi31.Spec, error) {
	paths, err := g.genRoutes(router)
	if err != nil {
		return nil, err
	}

	return &openapi31.Spec{
		Openapi: "3.1",
		Info: openapi31.Info{
			Title:       router.Summary(),
			Summary:     P(router.Summary()),
			Description: P(router.Description()),
		},
		JSONSchemaDialect: nil,
		Servers:           nil,
		Paths:             paths,
		Webhooks:          nil,
		Components:        nil,
		Security:          nil,
		Tags:              nil,
		ExternalDocs:      nil,
		MapOfAnything:     nil,
	}, nil
}

func (g *generatorImpl) genRoutes(router Router) (*openapi31.Paths, error) {
	paths := &openapi31.Paths{
		MapOfPathItemValues: make(map[string]openapi31.PathItem),
		MapOfAnything:       make(map[string]interface{}),
	}

	for _, route := range router.Routes() {
		err := g.genRoute(paths, "", route)
		if err != nil {
			return nil, err
		}
	}

	return paths, nil
}

func (g *generatorImpl) genRoute(paths *openapi31.Paths, prefix string, route chi.Route) error {
	pattern := path.Join(prefix, route.Pattern)

	item := openapi31.PathItem{
		Summary:       nil,
		Description:   nil,
		Servers:       nil,
		Parameters:    nil,
		MapOfAnything: make(map[string]interface{}),
	}

	// this routes handlers
	for method, handler := range route.Handlers {
		operation := &openapi31.Operation{
			ExternalDocs:  nil,
			ID:            nil,
			Tags:          nil,
			Summary:       nil,
			Description:   nil,
			Parameters:    nil,
			RequestBody:   nil,
			Responses:     nil,
			Callbacks:     nil,
			Deprecated:    nil,
			Security:      nil,
			Servers:       nil,
			MapOfAnything: nil,
		}

		if typed, ok := handler.(Handler); ok {
			info := typed.Info()

			if info.Hidden() {
				continue
			}

			operation.Summary = P(info.Summary())
			operation.Description = P(info.Description())

			request, err := g.requestBody(info.Input())
			if err != nil {
				return err
			}
			operation.RequestBody = request

			response, err := g.responseBody(info.Output())
			if err != nil {
				return err
			}
			operation.Responses = &openapi31.Responses{
				Default:                        response,
				MapOfResponseOrReferenceValues: g.defaultResponses,
				MapOfAnything:                  nil,
			}
		}

		switch strings.ToUpper(method) {
		case "GET":
			item.Get = operation
			break
		case "PUT":
			item.Put = operation
			break
		case "POST":
			item.Post = operation
			break
		case "DELETE":
			item.Delete = operation
			break
		case "OPTIONS":
			item.Options = operation
			break
		case "HEAD":
			item.Head = operation
			break
		case "PATCH":
			item.Patch = operation
			break
		case "TRACE":
			item.Trace = operation
			break
		case "CONNECT", "*":
			continue
		default:
			return errors.New("unknown method: " + method)
		}
	}

	paths.MapOfPathItemValues[pattern] = item

	// recursive sub routes
	if route.SubRoutes != nil {
		for _, subroute := range route.SubRoutes.Routes() {
			if err := g.genRoute(paths, pattern, subroute); err != nil {
				return err
			}
		}
	}

	return nil
}

func (gen *generatorImpl) requestBody(typ *HandlerType) (*openapi31.RequestBodyOrReference, error) {
	if typ == nil {
		return nil, nil
	}

	spec, err := gen.reflector.Reflect(reflect.New(typ.GoType).Interface())
	if err != nil {
		return nil, err
	}

	content := map[string]openapi31.MediaType{}
	for _, mime := range typ.HttpType {
		content[mime] = openapi31.MediaType{
			Schema: structToMap(spec),
		}
	}

	return &openapi31.RequestBodyOrReference{
		RequestBody: &openapi31.RequestBody{
			Description:   nil,
			Content:       content,
			Required:      nil,
			MapOfAnything: nil,
		},
	}, nil
}

func (gen *generatorImpl) responseBody(typ *HandlerType) (*openapi31.ResponseOrReference, error) {
	if typ == nil {
		return nil, nil
	}

	spec, err := gen.reflector.Reflect(reflect.New(typ.GoType).Interface())
	if err != nil {
		return nil, err
	}

	content := map[string]openapi31.MediaType{}
	for _, mime := range typ.HttpType {
		content[mime] = openapi31.MediaType{
			Schema: structToMap(spec),
		}
	}

	response := &openapi31.Response{
		Description:   "",
		Headers:       nil,
		Links:         nil,
		MapOfAnything: nil,

		Content: content,
	}

	return &openapi31.ResponseOrReference{
		Response: response,
	}, nil
}
