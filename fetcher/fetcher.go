package fetcher

import (
	"net/http"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)
type Middleware func(RoundTripperFunc) RoundTripperFunc
type MiddlewareTransport func(*http.Request) (*http.Response, error)

// Chain multiple middlewares into one,
// taken from go-kit
func Chain(outer Middleware, middlewares ...Middleware) Middleware {
	return func(next RoundTripperFunc) RoundTripperFunc {
		for i := len(middlewares) - 1; i >= 0; i-- { // reverse
			next = middlewares[i](next)
		}
		return outer(next)
	}
}

func NoopMiddleware(next RoundTripperFunc) RoundTripperFunc {
	return func(req *http.Request) (*http.Response, error) {
		return next(req)
	}
}

func (m MiddlewareTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}
