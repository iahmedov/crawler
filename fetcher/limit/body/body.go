package body

import (
	"io"
	"net/http"

	"github.com/iahmedov/crawler/fetcher"
	"github.com/pkg/errors"
)

func init() {
	fetcher.RegisterMiddlewareFactory("limit_body", New)
}

func New(config fetcher.Config) (fetcher.Middleware, error) {
	limitIfc, ok := config["limit"]
	if !ok {
		return nil, errors.New("limit not provided")
	}

	limit, ok := limitIfc.(int)
	if !ok {
		return nil, errors.New("limit should have int type")
	}

	return WithBodyLimit(int64(limit)), nil
}

type readCloser struct {
	io.Reader
	io.Closer
}

func (r *readCloser) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func (r *readCloser) Close() error {
	return r.Closer.Close()
}

func LimitedReadCloser(size int64, rc io.ReadCloser) io.ReadCloser {
	return &readCloser{
		Reader: io.LimitReader(rc, size),
		Closer: rc,
	}
}

func WithBodyLimit(size int64) fetcher.Middleware {
	return func(next fetcher.RoundTripperFunc) fetcher.RoundTripperFunc {
		return func(req *http.Request) (*http.Response, error) {
			response, err := next(req)
			if err != nil {
				return response, err
			}
			response.Body = LimitedReadCloser(size, response.Body)
			return response, err
		}
	}
}
