package test

import (
	"encoding/json"
	"github.com/benni-tec/gocart/gotrac"
	"net/http"
	"testing"
)

func TestDocs(t *testing.T) {
	gen := gotrac.NewGenerator()
	router := gotrac.Default().
		WithSummary("Fotobox Firmware (Backend)").
		WithDescription("some kind of description")

	router.MethodFunc(http.MethodGet, "/ping", pong).
		WithSummary("Ping -> Pong").
		WithDescription("Returns a pong").
		WithOutput(gotrac.Json[PongResponse]())

	docs, err := gen.Generate(router)
	if err != nil {
		t.Fatal(err)
	}

	js, err := json.MarshalIndent(docs, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	router.WithDocs("/openapi.json", gen)
	router.WithSwaggerUI("/swagger", "/openapi.json", "API Explorer", nil)

	t.Log(string(js))
}

func pong(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)

	_, _ = writer.Write([]byte("{\"pong\": true}"))
}

type PongResponse struct {
	Pong bool `json:"pong"`
}
