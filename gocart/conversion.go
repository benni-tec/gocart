package gocart

import (
	"github.com/benni-tec/gocart/gotrac"
	"github.com/benni-tec/gocart/middleware"
	"io"
	"net/http"
)

type Cart interface {
	gotrac.Handler
	WithInfo(fn func(info *CartInformation)) Cart
}

type CartFunc[TInput any, TOutput any] func(request *Request[TInput], writer HeaderWriter) (*TOutput, error)

type Converter[T any] interface {
	Serialize(body *T, headers http.Header) ([]byte, error)
	Deserialize(data []byte, headers http.Header) (*T, error)
	Type() *gotrac.HandlerType
}

type cartImpl[TInput any, TOutput any] struct {
	info   CartInformation
	input  Converter[TInput]
	output Converter[TOutput]

	handler CartFunc[TInput, TOutput]
}

func IO[TInput any, TOutput any](input Converter[TInput], output Converter[TOutput], h CartFunc[TInput, TOutput]) Cart {
	return &cartImpl[TInput, TOutput]{
		input:   input,
		output:  output,
		handler: h,

		info: CartInformation{
			summary:     "",
			description: "",
			hidden:      false,
		},
	}
}

func I[TInput any](input Converter[TInput], h CartFunc[TInput, any]) Cart {
	return IO(input, nil, h)
}

func O[TOutput any](output Converter[TOutput], h CartFunc[any, TOutput]) Cart {
	return IO(nil, output, h)
}

func A(h CartFunc[any, any]) Cart {
	return IO(nil, nil, h)
}

func (cart *cartImpl[TInput, TOutput]) Info() *gotrac.HandlerInformation {
	info := cart.info

	handler := &gotrac.HandlerInformation{
		Information: gotrac.Information{
			Summary:     info.summary,
			Description: info.description,
		},
		Hidden: info.hidden,
	}

	if cart.input != nil {
		handler.Input = cart.input.Type()
	}

	if cart.output != nil {
		handler.Output = cart.output.Type()
	}

	return handler
}

func (cart *cartImpl[TInput, TOutput]) WithInfo(fn func(info *CartInformation)) Cart {
	if fn != nil {
		fn(&cart.info)
	}

	return cart
}

func (cart *cartImpl[TInput, TOutput]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setTypes(w.Header(), "Accepts", cart.input.Type().HttpType)
	setTypes(w.Header(), "Content-Type", cart.output.Type().HttpType)

	errors := middleware.GetErrors(r.Context())

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errors.AddError(err)
		return
	}

	input, err := cart.input.Deserialize(body, r.Header)
	if err != nil {
		errors.AddError(err)
		return
	}

	output, err := cart.handler(ConvertRequest(r, input), w)
	if err != nil {
		errors.AddError(err)
		return
	}

	body, err = cart.output.Serialize(output, w.Header())
	if err != nil {
		errors.AddError(err)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		errors.AddError(err)
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

type CartInformation struct {
	summary     string
	description string
	hidden      bool
}

func (actor *CartInformation) WithSummary(summary string) *CartInformation {
	actor.summary = summary
	return actor
}

func (actor *CartInformation) WithDescription(description string) *CartInformation {
	actor.description = description
	return actor
}

func (actor *CartInformation) WithHidden(hidden bool) *CartInformation {
	actor.hidden = hidden
	return actor
}
