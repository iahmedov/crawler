package circuitbreaker

import (
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/iahmedov/crawler/fetcher"
	"github.com/sony/gobreaker"
)

func init() {
	fetcher.RegisterMiddlewareFactory("circuit_breaker", New)
}

type breakerConfig struct {
	MaxRequests uint32        `mapstructure:"max_requests"`
	Interval    time.Duration `mapstructure:"interval"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

func New(config fetcher.Config) (fetcher.Middleware, error) {
	var parsedConfig breakerConfig
	if err := mapstructure.Decode(config, &parsedConfig); err != nil {
		return nil, err
	}

	settings := gobreaker.Settings{
		Name:        "gobreaker",
		MaxRequests: parsedConfig.MaxRequests,
		Interval:    parsedConfig.Interval,
		Timeout:     parsedConfig.Timeout,
	}
	return CircuitBreaker(settings), nil
}

func CircuitBreaker(settings gobreaker.Settings) fetcher.Middleware {
	cb := gobreaker.NewCircuitBreaker(settings)
	return func(next fetcher.RoundTripperFunc) fetcher.RoundTripperFunc {
		return func(req *http.Request) (*http.Response, error) {
			response, err := cb.Execute(func() (interface{}, error) { return next(req) })
			return response.(*http.Response), err
		}
	}
}
