package test

import (
	"github.com/benni-tec/gocart/gocart"
	"github.com/benni-tec/gocart/gotrac"
	"net/http"
	"testing"
)

func TestConversion(t *testing.T) {
	router := gotrac.Default()

	router.Method(http.MethodGet, "/ping", gocart.IO(gocart.Json[PingRequest](), gocart.Json[PongResponse](), ping))
	router.Method(http.MethodGet, "/file", gocart.O(gocart.Binary(), download))

	// +++ Test +++
	// TODO
}

func download(_ *gocart.Request[any], _ gocart.HeaderWriter) (*[]byte, error) {
	return nil, nil
}

func ping(request *gocart.Request[PingRequest], _ gocart.HeaderWriter) (*PongResponse, error) {
	parity := request.Body().N % 2

	return &PongResponse{
		Pong: parity == 0,
	}, nil
}

type PingRequest struct {
	N int `json:"n"`
}
