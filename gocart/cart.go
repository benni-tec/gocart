package gocart

import (
	"github.com/benni-tec/gocart/goflag"
	"github.com/benni-tec/gocart/gotrac"
	"github.com/benni-tec/gocart/middleware"
	"io"
	"net/http"
	"reflect"
)

// Cart represents a gotrac.EndpointFlag that automatically can (de)serialize the request/response
type Cart interface {
	goflag.EndpointFlag
	WithInfo(fn func(info *CartInformation)) Cart
}

// CartFunc is the actual handler that gets the deserialized request (and a HeaderWriter)
// and returns the response body that will then be serialized by the cart.
type CartFunc[TInput any, TOutput any] func(request *Request[TInput], writer HeaderWriter) (*TOutput, error)

type cartImpl[TInput any, TOutput any] struct {
	info   CartInformation
	input  Serializer[TInput]
	output Serializer[TOutput]

	handler CartFunc[TInput, TOutput]
}

// IO is used to define a Cart that deserializes it's input and serializes the output.
func IO[TInput any, TOutput any](input Serializer[TInput], output Serializer[TOutput], h CartFunc[TInput, TOutput]) Cart {
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

// I is used to define a Cart that deserializes it's input but does NOT serialize the output.
//
// However, TOutput can still be specified to allow for meta-fields!
func I[TInput any, TOutput any](input Serializer[TInput], h CartFunc[TInput, TOutput]) Cart {
	return IO[TInput, TOutput](input, nil, h)
}

// O is used to define a Cart that does NOT deserialize it's input but does serialize the output.
//
// However, TInput can still be specified to allow for meta-fields!
func O[TInput any, TOutput any](output Serializer[TOutput], h CartFunc[TInput, TOutput]) Cart {
	return IO[TInput, TOutput](nil, output, h)
}

// A is used to define a Cart that does NOT deserialize it's input and does NOT serialize the output.
//
// However, TInput and TOutput can still be specified to allow for meta-fields!
func A[TInput any, TOutput any](h CartFunc[TInput, TOutput]) Cart {
	return IO[TInput, TOutput](nil, nil, h)
}

func (cart *cartImpl[TInput, TOutput]) Info() *goflag.EndpointInformation {
	info := cart.info

	handler := &goflag.EndpointInformation{
		Information: goflag.Information{
			Summary:     info.summary,
			Description: info.description,
		},
		Hidden: info.hidden,
	}

	if cart.input != nil {
		handler.Input = cart.input.Type()
	} else {
		handler.Input = gotrac.None[TInput]()
	}

	if cart.output != nil {
		handler.Output = cart.output.Type()
	} else {
		handler.Output = gotrac.None[TOutput]()
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
	setHeader(w.Header(), "Accepts", cart.input)
	setHeader(w.Header(), "Content-Type", cart.output)

	errors := middleware.GetErrors(r.Context())

	input, err := cart.decode(r)
	if err != nil {
		errors.AddError(err)
		return
	}

	output, err := cart.handler(wrapToBodyRequest[TInput](r, input), w)
	if err != nil {
		errors.AddError(err)
		return
	}

	err = cart.encode(w, output)
	if err != nil {
		errors.AddError(err)
		return
	}
}

// +++ Codec +++

var encoderFactories = []EncoderFactory{
	NewHeaderEncoder,
	// TODO: add support for cookies
}

var decoderFactories = []DecoderFactory{
	NewHeaderDecoder,
	NewFormDecorder,
	NewPathDecoder,
	NewQueryDecoder,
	// TODO: add support for cookies
}

func (cart *cartImpl[TInput, TOutput]) decode(r *http.Request) (*TInput, error) {
	var input *TInput
	if cart.input == nil {
		input = new(TInput)
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		input, err = cart.input.Deserialize(body, r.Header)
		if err != nil {
			return nil, err
		}
	}

	// populate path, meta, formData and query parameters
	var decoders []Decoder
	for _, factory := range decoderFactories {
		decoders = append(decoders, factory(r))
	}

	val := reflect.Indirect(reflect.ValueOf(input))
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return input, nil
	}

	for i := range typ.NumField() {
		structField := typ.Field(i)
		field := val.Field(i)

		for _, dec := range decoders {
			vals, err := dec.Decode(structField)
			if err != nil {
				return nil, err
			}

			err = AssignPrimitives(field, vals)
			if err != nil {
				return nil, err
			}
		}

		// TODO: add support for json tags, check constraints, ...
		// https://pkg.go.dev/github.com/swaggest/jsonschema-go#Reflector.Reflect
	}

	return input, nil
}

func (cart *cartImpl[TInput, TOutput]) encode(w http.ResponseWriter, output *TOutput) error {
	var encoders []Encoder
	for _, factory := range encoderFactories {
		encoders = append(encoders, factory(w))
	}

	if output != nil {
		val := reflect.Indirect(reflect.ValueOf(output))
		if val.IsZero() && val.Type().Kind() == reflect.Struct {
			typ := val.Type()

			for i := range typ.NumField() {
				structField := typ.Field(i)
				field := val.Field(i)

				for _, dec := range encoders {
					err := dec.Encode(field, structField)
					if err != nil {
						return err
					}
				}

				// TODO: add support for json tags, check constraints, ...
				// https://pkg.go.dev/github.com/swaggest/jsonschema-go#Reflector.Reflect
			}
		}
	}

	if cart.output != nil {
		body, err := cart.output.Serialize(output, w.Header())
		if err != nil {
			return err
		}

		_, err = w.Write(body)
		if err != nil {
			return err
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	return nil
}
