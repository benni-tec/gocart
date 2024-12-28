package gocrew

import (
	"encoding/json"
	"errors"
	"github.com/benni-tec/gocart/gotrac"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go/openapi31"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"net/http"
	"path"
	"reflect"
	"strings"
)

// +++ Spec +++

type OpenApi31Spec openapi31.Spec

func (spec *OpenApi31Spec) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	j, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(j)
}

func (spec *OpenApi31Spec) WithUI(title string, basePath string, docPattern string, config *swg.Config) http.Handler {
	if config != nil {
		ui := swgui.NewWithConfig(*config)
		return ui(title, docPattern, basePath)
	} else {
		return swgui.New(title, docPattern, basePath)
	}
}

// +++ Generator +++

type openapi31Generator struct {
	reflector        *openapi31.Reflector
	defaultResponses map[string]openapi31.ResponseOrReference
}

func (g *openapi31Generator) Generate(router gotrac.Router) (*OpenApi31Spec, error) {
	paths, err := g.genRoutes(router)
	if err != nil {
		return nil, err
	}

	info := router.Info()
	return &OpenApi31Spec{
		Openapi: "3.1.0",
		Info: openapi31.Info{
			Title:       info.Summary,
			Summary:     gotrac.P(info.Summary),
			Description: gotrac.P(info.Description),
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

func (g *openapi31Generator) genRoutes(router gotrac.Router) (*openapi31.Paths, error) {
	operations := make(map[string]map[string]*openapi31.Operation)
	err := chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		typed, isTyped := handler.(gotrac.Handler)
		if !isTyped {
			return nil
		}

		info := typed.Info()
		if info.Hidden {
			return nil
		}

		operation := &openapi31.Operation{
			Summary:     gotrac.P(info.Summary),
			Description: gotrac.P(info.Description),
		}

		request, err := g.requestBody(info.Input)
		if err != nil {
			return err
		}
		operation.RequestBody = request

		response, err := g.responseBody(info.Output)
		if err != nil {
			return err
		}
		operation.Responses = &openapi31.Responses{
			Default:                        response,
			MapOfResponseOrReferenceValues: g.defaultResponses,
			MapOfAnything:                  nil,
		}

		// ensure map has key
		_, ok := operations[route]
		if !ok {
			operations[route] = make(map[string]*openapi31.Operation)
		}

		operations[route][method] = operation
		return nil
	})

	if err != nil {
		return nil, err
	}

	paths := &openapi31.Paths{
		MapOfPathItemValues: make(map[string]openapi31.PathItem),
		MapOfAnything:       make(map[string]interface{}),
	}

	for route, x := range operations {
		item := openapi31.PathItem{}

		for method, operation := range x {
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
				return nil, errors.New("unknown method: " + method)
			}
		}

		paths.MapOfPathItemValues[route] = item
	}

	return paths, nil
}

func (g *openapi31Generator) genRoute(paths *openapi31.Paths, prefix string, route chi.Route) error {
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

		if typed, ok := handler.(gotrac.Handler); ok {
			info := typed.Info()

			if info.Hidden {
				continue
			}

			operation.Summary = gotrac.P(info.Summary)
			operation.Description = gotrac.P(info.Description)

			request, err := g.requestBody(info.Input)
			if err != nil {
				return err
			}
			operation.RequestBody = request

			response, err := g.responseBody(info.Output)
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

func (gen *openapi31Generator) requestBody(typ *gotrac.HandlerType) (*openapi31.RequestBodyOrReference, error) {
	if typ == nil {
		return nil, nil
	}

	spec, err := gen.reflector.Reflect(reflect.New(typ.GoType).Interface())
	if err != nil {
		return nil, err
	}

	schema, err := structToMap(spec)
	if err != nil {
		return nil, err
	}

	content := map[string]openapi31.MediaType{}
	for _, mime := range typ.HttpType {
		content[mime] = openapi31.MediaType{
			Schema: schema,
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

func (gen *openapi31Generator) responseBody(typ *gotrac.HandlerType) (*openapi31.ResponseOrReference, error) {
	if typ == nil {
		return nil, nil
	}

	spec, err := gen.reflector.Reflect(reflect.New(typ.GoType).Interface())
	if err != nil {
		return nil, err
	}

	schema, err := structToMap(spec)
	if err != nil {
		return nil, err
	}

	content := map[string]openapi31.MediaType{}
	for _, mime := range typ.HttpType {
		content[mime] = openapi31.MediaType{
			Schema: schema,
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
