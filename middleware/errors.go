package middleware

import (
	"context"
	"github.com/benni-tec/gocart/utils"
	"github.com/go-chi/chi/v5/middleware"
)

var errorCache = utils.NewCache[[]error]()

type Errors interface {
	Id() string
	Errors() []error
	AddError(err error)
}

func GetErrors(ctx context.Context) Errors {
	withId := middleware.GetReqID(ctx)
	return &errorsImpl{id: withId}
}

type errorsImpl struct {
	id string
}

func (e *errorsImpl) Id() string {
	return e.id
}

func (c *errorsImpl) Errors() []error {
	return errorCache.Get(c.Id())
}

func (c *errorsImpl) AddError(err error) {
	errs := c.Errors()
	errs = append(errs, err)
	errorCache.Set(c.Id(), errs)
}
