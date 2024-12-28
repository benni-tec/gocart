package middleware

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// ErrorMiddleware allows for Errors to be attached to a request.
// It uses a global cache with the request id as a key, it therefore requires chi`s RequestId middleware!
//
// If errors are present a 500 code will be returned and the errors encoded as json.
//
// The GetErrors function and the Errors interface can be used without this middleware (only requiring chi`s RequestId)
// so you can write your own error handler!
func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		errors := GetErrors(r.Context()).Errors()
		if errors == nil || len(errors) == 0 {
			return
		}

		errs := make([]string, len(errors))
		for _, err := range errors {
			errs = append(errs, err.Error())
		}

		js, err := json.Marshal(struct {
			Errors []string `json:"errors"`
		}{
			Errors: errs,
		})

		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(js)
		if err != nil {
			panic(err)
		}
	})
}

// Errors can be used to read and add errors to the error cache for this request!
//
// This requires only chi`s RequestId middleware, while you can use the ErrorMiddleware to handle the returned errors
// you can also write your own!
type Errors interface {
	// Id returns the request id given to the context`s request by chi`s RequestId middlware
	Id() string
	// Errors returns the attached errors
	Errors() []error
	// AddError can be used to add an error
	AddError(err error)

	// Done should only be called once the response has reached to error handler,
	// and the errors can be deleted from the cache.
	//
	// Since this clears the cache for the associated request id,
	// using Errors or AddError after this call will cause a panic!
	Done()
}

// GetErrors retrieves the Errors interface for a give request`s context.
//
// This requires only chi`s RequestId middleware, while you can use the ErrorMiddleware to handle the returned errors
// you can also write your own!
func GetErrors(ctx context.Context) Errors {
	withId := middleware.GetReqID(ctx)
	if withId == "" {
		panic("middleware: no request id found! Try adding the RequestId middleware.")
	}

	return &errorsImpl{id: withId}
}

// +++ Cache +++

var errorCache = map[string][]error{}

type errorsImpl struct {
	id string
}

func (e *errorsImpl) Id() string {
	return e.id
}

func (c *errorsImpl) Errors() []error {
	return errorCache[c.id]
}

func (c *errorsImpl) AddError(err error) {
	errs := c.Errors()
	errs = append(errs, err)
	errorCache[c.id] = errs
}

func (c *errorsImpl) Done() {
	delete(errorCache, c.id)
}
