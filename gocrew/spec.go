package gocrew

import (
	"encoding/json"
	"github.com/swaggest/openapi-go/openapi31"
	swg "github.com/swaggest/swgui"
	swgui "github.com/swaggest/swgui/v5emb"
	"net/http"
)

type Spec openapi31.Spec

func (spec *Spec) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	j, err := json.Marshal(spec)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(j)
}

func (spec *Spec) WithUI(title string, basePath string, docPattern string, config *swg.Config) http.Handler {
	if config != nil {
		ui := swgui.NewWithConfig(*config)
		return ui(title, docPattern, basePath)
	} else {
		return swgui.New(title, docPattern, basePath)
	}
}
