package goflag

import (
	"net/http"
	"reflect"
)

// EndpointFlag is just a http.Handler which can also provide EndpointInformation
type EndpointFlag interface {
	http.Handler
	Info() *EndpointInformation
}

// Type defines both the type used in go and the http content type (MIME-type)
type Type struct {
	GoType   reflect.Type
	HttpType []string
}

// EndpointInformation contains the information that can be set for a handler.
// This is only readable since handler can be anything provided to gotrac.
// Once the handler is registered with a Router a Route is returned where the information can be edited.
type EndpointInformation struct {
	Information
	Input  *Type
	Output *Type
	Hidden bool
}

func (c *EndpointInformation) WithSummary(summary string) *EndpointInformation {
	c.Information.WithSummary(summary)
	return c
}

func (c *EndpointInformation) WithDescription(description string) *EndpointInformation {
	c.Information.WithDescription(description)
	return c
}

func (c *EndpointInformation) WithInput(typ *Type) *EndpointInformation {
	c.Input = typ
	return c
}

func (c *EndpointInformation) WithOutput(typ *Type) *EndpointInformation {
	c.Output = typ
	return c
}

func (c *EndpointInformation) WithHidden(hidden bool) *EndpointInformation {
	c.Hidden = hidden
	return c
}

type flaggedEndpoint struct {
	http.Handler
	flag[EndpointInformation]
}

func FlagEndpoint(endpoint http.Handler, info ...EndpointInformation) EndpointFlag {
	var _info EndpointInformation
	if len(info) > 0 {
		_info = info[0]
	}

	return &flaggedEndpoint{
		Handler: endpoint,
		flag:    flag[EndpointInformation]{info: _info},
	}
}
