package gotrac

import (
	"encoding/json"
	"github.com/swaggest/openapi-go"
	"net/http"
)

func specHandler(s openapi.SpecSchema) http.Handler {
	j, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(j)
	})
}
