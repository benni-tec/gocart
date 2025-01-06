package gocrew

import (
	"encoding/json"
	"github.com/benni-tec/gocart/goflag"
	"github.com/go-chi/chi/v5"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"net/http"
	"reflect"
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
	title            string
	defaultTag       openapi31.Tag
	defaultResponses map[string]openapi31.ResponseOrReference
}

func (gen *openapi31Generator) Generate(router chi.Routes) (*OpenApi31Spec, error) {
	reflector := openapi31.NewReflector()
	reflector.Spec = &openapi31.Spec{Openapi: "3.1.0"}

	if r, ok := router.(goflag.InformationFlag); ok {
		reflector.Spec.Info.
			WithSummary(r.Info().Summary).
			WithDescription(r.Info().Description)
	}

	reflector.Spec.Info.WithTitle(gen.title)

	hasDefaultTag := false

	err := Walk(
		router,
		func(method string, route string, handler http.Handler, controller goflag.ControllerFlag) error {
			ctx, err := reflector.NewOperationContext(method, route)
			if err != nil {
				return err
			}

			typed, isTyped := handler.(goflag.EndpointFlag)
			if !isTyped {
				return nil
			}

			info := typed.Info()
			if info.Hidden {
				return nil
			}

			var tag string
			if controller != nil {
				tag = controller.Info().Name
			} else {
				hasDefaultTag = true
				tag = gen.defaultTag.Name
			}

			ctx.SetSummary(info.Summary)
			ctx.SetDescription(info.Description)
			ctx.SetTags(tag)

			// schemas
			// TODO: set proper contentType
			if info.Input != nil {
				dummy := reflect.New(info.Input.GoType).Interface()

				if len(info.Input.HttpType) == 0 {
					ctx.AddReqStructure(dummy, openapi.WithHTTPStatus(http.StatusNoContent))
				}

				for _, typ := range info.Input.HttpType {
					ctx.AddReqStructure(dummy, openapi.WithContentType(typ))
				}
			}

			if info.Output != nil {
				dummy := reflect.New(info.Output.GoType).Interface()

				if len(info.Output.HttpType) == 0 {
					ctx.AddRespStructure(dummy, openapi.WithHTTPStatus(http.StatusNoContent))
				}

				for _, typ := range info.Output.HttpType {
					ctx.AddRespStructure(dummy, openapi.WithContentType(typ))
				}
			}

			return reflector.AddOperation(ctx)
		},
		func(controller goflag.ControllerFlag) error {
			info := controller.Info()

			tag := openapi31.Tag{}
			tag.WithName(info.Name)
			tag.WithDescription(info.Description)

			_ = append(reflector.Spec.Tags, tag)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	if hasDefaultTag {
		_ = prepend(reflector.Spec.Tags, gen.defaultTag)
	}

	spec := OpenApi31Spec(*reflector.Spec)
	return &spec, nil
}
