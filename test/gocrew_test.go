package test

import (
	"encoding/json"
	"github.com/benni-tec/gocart/gocrew"
	"github.com/benni-tec/gocart/gotrac"
	"net/http"
	"testing"
)

func TestDocs(t *testing.T) {
	gen := gocrew.OpenApi31()
	router := gotrac.Default().WithInfo(func(info *gotrac.RouterInformation) {
		info.WithSummary("Fotobox Firmware (Backend)").
			WithDescription("some kind of description")
	})

	router.MethodFunc(http.MethodGet, "/ping", pong).WithInfo(func(info *gotrac.RouteInformation) {
		info.WithSummary("Ping -> Pong").
			WithDescription("Returns a pong").
			WithOutput(gotrac.Json[PongResponse]())
	})

	docs, err := gen.Generate(router)
	if err != nil {
		t.Fatal(err)
	}

	js, err := json.MarshalIndent(docs, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	spec, err := gen.Generate(router)
	if err != nil {
		t.Fatal(err)
	}

	router.Mount("/openapi.json", spec)
	router.Mount("/swagger", spec.WithUI("API Explorer", "/swagger", "/openapi.json", nil))

	t.Log(string(js))
}

func pong(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)

	_, _ = writer.Write([]byte("{\"pong\": true}"))
}

type PongResponse struct {
	Pong bool `json:"pong"`
}
