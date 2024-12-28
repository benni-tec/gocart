package gocart

import "net/http"

// HeaderWriter allows a CartFunc to write and read to the response header without being able to write the body.
// This helps to ensure the proper write flow, i.e. headers can not be written after the body.
type HeaderWriter interface {
	// Header returns the meta map that will be sent by
	// [ResponseWriter.WriteHeader]. The [Header] map also is the mechanism with which
	// [Handler] implementations can set HTTP trailers.
	//
	// Changing the meta map after a call to [ResponseWriter.WriteHeader] (or
	// [ResponseWriter.Write]) has no effect unless the HTTP status code was of the
	// 1xx class or the modified headers are trailers.
	//
	// There are two ways to set Trailers. The preferred way is to
	// predeclare in the headers which trailers you will later
	// send by setting the "Trailer" meta to the names of the
	// trailer keys which will come later. In This case, those
	// keys of the Header map are treated as if they were
	// trailers. See the example. The second way, for trailer
	// keys not known to the [Handler] until after the first [ResponseWriter.Write],
	// is to prefix the [Header] map keys with the [TrailerPrefix]
	// constant value.
	//
	// To suppress automatic response headers (such as "Date"), set
	// their value to nil.
	Header() http.Header

	// WriteHeader sends an HTTP response meta with the provided
	// status code.
	//
	// If WriteHeader is not called explicitly, the first call to Write
	// will trigger an implicit WriteHeader(gocart.StatusOK).
	// Thus explicit calls to WriteHeader are mainly used to
	// send error codes or 1xx informational responses.
	//
	// The provided code must be a valid HTTP 1xx-5xx status code.
	// Any number of 1xx headers may be written, followed by at most
	// one 2xx-5xx meta. 1xx headers are sent immediately, but 2xx-5xx
	// headers may be buffered. Use the Flusher interface to send
	// buffered data. The meta map is cleared when 2xx-5xx headers are
	// sent, but not with 1xx headers.
	//
	// The server will automatically send a 100 (Continue) meta
	// on the first read from the request body if the request has
	// an "Expect: 100-continue" meta.
	WriteHeader(statusCode int)
}
