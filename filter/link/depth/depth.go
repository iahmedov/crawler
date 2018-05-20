package depth

import (
	"net/url"

	"github.com/iahmedov/crawler/filter"
	filterlink "github.com/iahmedov/crawler/filter/link"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	filterlink.RegisterFilterFactory("depth", New)
}

type depthConfig struct {
	Depth uint32 `mapstructure:"include"`
}

func New(config filter.FilterConfig) (filterlink.Filter, error) {
	c := &depthConfig{}
	if err := mapstructure.Decode(config, c); err != nil {
		return nil, errors.Wrap(err, "could not parse depth link filter config")
	}

	return Depth(c.Depth), nil
}

func Depth(depth uint32) filterlink.Filter {
	return func(u url.URL, urlDepth uint32) filter.State {
		if urlDepth < depth {
			return filter.StatePositive
		}
		return filter.StateNegative
	}
}
