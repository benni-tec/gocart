# <img src=".github/gokart.png" style="height: 1em; position: relative; top: 0.15em"> gocart

This project is based upon the [chi](https://github.com/go-chi/chi) router 
and is strongly inspired by the [swaggest](https://github.com/swaggest/rest) project.

## Disclaimer
This project is currently under development and, while I'm using it in my projects, 
I would not recommend using this in production-level software.

This was a holiday project and is thus not complete at all. I will add functionality when I need it, 
however some more complex features (like cookies) are still missing.

If you are interested in using this framework feel free to open PRs or DM me!

## Motivation
I created this mostly to re-learn golang after spending an extended time working with [ABP](https://abp.io/) 
(Angular + C# with MVC and [nswag](https://github.com/RicoSuter/NSwag)). 
My goal was to find something similar to ASP.NET to autogenerate documentation, preferably as a OpenAPI specification.

While [swaggest](https://github.com/swaggest/rest) does what I need, I was missing some functionality that was key for me:
- Support for files both in requests and responses,
- Declaring types without actually deserializing the body,
- Support for serialization other than JSON, like Yaml or XML.

## Concept
This project is split into 3 and a half parts:
- `gotrac` - Router with meta-data information
- `gocart` - Automatic (de)serialization for handlers
- `gocrew` - Generate documentation

Both `gocart` and `gocrew` build upon `gotrac`, but not on each other! 
`gotrac` contains all the information needed to generate documentation, `gocart.Cart` only implements the `gotrac.Handler` 
with automatic (de)serialization but does not add more information for documentation that could not also be given directly to `gotrac`.

## <img src=".github/race-track.png" style="height: 1em; position: relative; top: 0.15em"> `gotrac`
This package provides a fully `net/http` compatible router, i.e. the `Router` itself is a `http.Handler`
and it follows the middleware standard.

The API is similar to `chi.Router` however it is not compatible, requiring a `gotrac.Handler` instead of a simple `http.Handler`.
The main difference is, that the `gotrac.Handler` can provide a `HandlerInformation` struct which is the used to construct the `gotrac.Route`.

While the `gotrac.Handler` is just an interface that you can implement arbitrarily, 
`gotrac.Route` represents a handler that was registered to a `gotrac.Router` and has editable information.
Basically it only wraps the actual `http.Handler` and copies the information.

We also provide a fluent-ish API to modify/set information:

```go
package main

import (
	"github.com/benni-tec/gocart/gotrac"
	"net/http"
)

func main() {
	router := gotrac.Default().WithInfo(func(info *gotrac.RouterInformation) {
		info.WithSummary("Fotobox Firmware (Backend)").
			WithDescription("some kind of description")
	})

	router.MethodFunc(http.MethodGet, "/ping", pong).
		WithInfo(func(info *gotrac.RouteInformation) {
            info.WithSummary("Ping -> Pong"). 
                 WithDescription("Returns a pong").
                 WithOutput(gotrac.Json[PongResponse]())
			
			// --- OR ---
			info.Summary = "Ping -> Pong"
			info.Description = "Returns a pong"
			info.Output = gotrac.Json[PongResponse]()
        })
	
	http.ListenAndServe(":8080", router)
}

func pong(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)

	_, _ = writer.Write([]byte("{\"pong\": true}"))
}

type PongResponse struct {
	Pong bool `json:"pong"`
}
```

TODO: Explain json-schema and meta-data attributes.

## <img src=".github/crew.png" style="height: 1em; position: relative; top: 0.15em"> `gocrew`
This package contains Generators for generating documentation from a `gotrac.Router`.

Notably this currently only includes the `gocrew.OpenApi31()` generator which generates a OpenAPI 3.1.0 compliant specification.
This `gocrew.OpenApi31Spec` also implements `http.Handler` and can thus be used to serve the generated specification.
It can also create a Swagger UI `http.Handler` to serve nice GUI.

```go
package main

import (
	"github.com/benni-tec/gocart/gocrew"
	"github.com/benni-tec/gocart/gotrac"
	"log"
	"net/http"
)

func main() {
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
	
	spec, err := gen.Generate(router)
	if err != nil {
		log.Fatal(err)
	}

	router.Mount("/openapi.json", spec)
	router.Mount("/swagger", spec.WithUI("API Explorer", "/swagger", "/openapi.json", nil))
}

func pong(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)

	_, _ = writer.Write([]byte("{\"pong\": true}"))
}

type PongResponse struct {
	Pong bool `json:"pong"`
}

```

## <img src=".github/gokart.png" style="height: 1em; position: relative; top: 0.15em"> `gocart`
This package builds upon `gotrac` with its main `gocart.Cart` type.
A `Cart` is a `gotrac.Handler` that implements (de)serialization of the body for you.

For this it uses an input and an output `gocart.Serializer` with implementations for json, yaml and xml
as well as a binary passthrough being provided here.

When handling a request it distinguishes between serializing the body with a `gocart.Serializer`
and decoding/encoding header information using a `gocart.Encoder` and a `gocart.Decoder`.

TODO: examples

## Attribution
- <a href="https://www.flaticon.com/free-icons/race-track" title="race track icons">Race track icons created by Freepik - Flaticon</a>
- <a href="https://www.flaticon.com/free-icons/go-kart" title="go kart icons">Go kart icons created by Leremy - Flaticon</a>
- <a href="https://www.flaticon.com/free-icons/studio" title="studio icons">Studio icons created by Leremy - Flaticon</a>