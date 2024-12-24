package gocart

import (
	"gocart/gotrac"
	"gocart/middleware"
	"io"
	"net/http"
)

type HandlerWithConversion interface {
	gotrac.Handler
	WithSummary(summary string) HandlerWithConversion
	WithDescription(description string) HandlerWithConversion
	WithHidden(hidden bool) HandlerWithConversion
}

type HandlerFunc[TInput any, TOutput any] func(request *Request[TInput], writer HeaderWriter) (*TOutput, error)

type Converter[T any] interface {
	Serialize(body *T, headers http.Header) ([]byte, error)
	Deserialize(data []byte, headers http.Header) (*T, error)
	Type() *gotrac.HandlerType
}

type conversionHandler[TInput any, TOutput any] struct {
	input  Converter[TInput]
	output Converter[TOutput]

	handler HandlerFunc[TInput, TOutput]

	summary     string
	description string
	hidden      bool
}

func WithConversion[TInput any, TOutput any](input Converter[TInput], output Converter[TOutput], h HandlerFunc[TInput, TOutput]) HandlerWithConversion {
	return &conversionHandler[TInput, TOutput]{
		input:   input,
		output:  output,
		handler: h,

		summary:     "",
		description: "",
		hidden:      false,
	}
}

func (actor *conversionHandler[TInput, TOutput]) Info() gotrac.HandlerInformation {
	return actor
}

func (actor *conversionHandler[TInput, TOutput]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setTypes(w.Header(), "Accepts", actor.input.Type().HttpType)
	setTypes(w.Header(), "Content-Type", actor.output.Type().HttpType)

	ctx := middleware.ContextWithErrors(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.AddError(err)
		return
	}

	input, err := actor.input.Deserialize(body, r.Header)
	if err != nil {
		ctx.AddError(err)
		return
	}

	output, err := actor.handler(ConvertRequest(r, input), w)
	if err != nil {
		ctx.AddError(err)
		return
	}

	body, err = actor.output.Serialize(output, w.Header())
	if err != nil {
		ctx.AddError(err)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		ctx.AddError(err)
		return
	}
}

func setTypes(header http.Header, key string, values []string) {
	header.Del(key)

	for _, value := range values {
		header.Add(key, value)
	}
}

// +++ Information +++

func (actor *conversionHandler[TInput, TOutput]) Summary() string {
	return actor.summary
}

func (actor *conversionHandler[TInput, TOutput]) Description() string {
	return actor.description
}

func (actor *conversionHandler[TInput, TOutput]) Input() *gotrac.HandlerType {
	return actor.input.Type()
}

func (actor *conversionHandler[TInput, TOutput]) Output() *gotrac.HandlerType {
	return actor.output.Type()
}

func (actor *conversionHandler[TInput, TOutput]) Hidden() bool {
	return actor.hidden
}

func (actor *conversionHandler[TInput, TOutput]) WithSummary(summary string) HandlerWithConversion {
	actor.summary = summary
	return actor
}

func (actor *conversionHandler[TInput, TOutput]) WithDescription(description string) HandlerWithConversion {
	actor.description = description
	return actor
}

func (actor *conversionHandler[TInput, TOutput]) WithHidden(hidden bool) HandlerWithConversion {
	actor.hidden = hidden
	return actor
}
