package test

import (
	"github.com/benni-tec/gocart/gocart"
	"github.com/benni-tec/gocart/gotrac"
	"net/http"
	"testing"
)

func TestConversion(t *testing.T) {
	router := gotrac.Default()

	router.Method(http.MethodGet, "/ping", gocart.WithConversion(gocart.Json[PingRequest](), gocart.Json[PongResponse](), ping))

	// +++ Test +++
	// TODO
}

func ping(request *gocart.Request[PingRequest], writer gocart.HeaderWriter) (*PongResponse, error) {

	return &PongResponse{
		Pong: request.Body().N%2 == 0,
	}, nil
}

type PingRequest struct {
	N int `json:"n"`
}
