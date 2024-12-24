package middleware

import (
	"context"
	"gocart/utils"
)

var errorCache = utils.NewCache[[]error]()

type Errors interface {
	ContextWithId
	Errors() []error
	AddError(err error)
}

func ContextWithErrors(ctx context.Context) Errors {
	withId := Context(ctx)
	return WithErrors(withId)
}

func WithErrors(ctx ContextWithId) Errors {
	return &errorsImpl{ContextWithId: ctx}
}

type errorsImpl struct {
	ContextWithId
}

func (c *errorsImpl) Errors() []error {
	return errorCache.Get(c.Id())
}

func (c *errorsImpl) AddError(err error) {
	errs := c.Errors()
	errs = append(errs, err)
	errorCache.Set(c.Id(), errs)
}
