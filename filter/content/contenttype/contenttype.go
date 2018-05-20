package porn

import (
	"github.com/iahmedov/crawler/filter"
	"github.com/iahmedov/crawler/filter/content"
	"github.com/iahmedov/crawler/model"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	content.RegisterFilterFactory("content_type", New)
}

type contentConfig struct {
	Keywords []string `mapstructure:"keywords"`
	Include  bool     `mapstructure:"include"`
}

func New(config filter.FilterConfig) (content.Filter, error) {
	c := &contentConfig{}
	if err := mapstructure.Decode(config, c); err != nil {
		return nil, errors.Wrap(err, "could not parse content type config")
	}

	return ContentTypeFilter(c.Keywords, c.Include), nil
}

func ContentTypeFilter(types []string, include bool) content.Filter {
	stateOnFound := filter.StateNeutral
	if include {
		stateOnFound = filter.StatePositive
	} else {
		stateOnFound = filter.StateNegative
	}

	return func(entry model.CrawlEntry) filter.State {
		contentType := entry.Header.Get("Content-type")
		for _, t := range types {
			if contentType == t {
				return stateOnFound
			}
		}
		return filter.StateNeutral
	}
}
