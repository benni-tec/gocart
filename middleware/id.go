package middleware

import (
	"context"
	"github.com/google/uuid"
	"gocart/utils"
	"net/http"
)

const idKey = "gocart:id"

var idSet = utils.NewSet[string]()

type ContextWithId interface {
	context.Context
	Id() string
}

func IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := newId()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), idKey, id)))
	})
}

func Context(ctx context.Context) ContextWithId {
	id := ctx.Value(idKey)
	if id == nil {
		panic("this contextImpl does not have an ID, maybe you need to add the IdMiddleware at startup")
	}

	str, ok := id.(string)
	if !ok {
		panic("this contextImpl does not have an ID, maybe you need to add the IdMiddleware at startup")
	}

	return &contextImpl{
		Context: ctx,
		id:      str,
	}
}

func newId() string {
	idSet.Lock()
	defer idSet.Unlock()

	// Find an id, that is net already used
	var id = ""
	for id == "" || idSet.Has(id) {
		id = uuid.New().String()
	}

	idSet.Add(id)
	return id
}

type contextImpl struct {
	context.Context
	id string
}

func (c *contextImpl) Id() string {
	return c.id
}
