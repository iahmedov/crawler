package domain

import (
	"net/url"
	"strings"

	"github.com/iahmedov/crawler/filter"
	filterlink "github.com/iahmedov/crawler/filter/link"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	filterlink.RegisterFilterFactory("domain", New)
}

type domainConfig struct {
	Domains []string `mapstructure:"domains"`
	Include bool     `mapstructure:"include"`
}

func New(config filter.FilterConfig) (filterlink.Filter, error) {
	c := &domainConfig{}
	if err := mapstructure.Decode(config, c); err != nil {
		return nil, errors.Wrap(err, "could not parse domain link filter config")
	}

	return Domain(c.Domains, c.Include), nil
}

func Domain(domains []string, include bool) filterlink.Filter {
	stateOnFound := filter.StateNeutral
	if include {
		stateOnFound = filter.StatePositive
	} else {
		stateOnFound = filter.StateNegative
	}

	return func(u url.URL, depth uint32) filter.State {
		host := u.Hostname()
		host = strings.ToLower(host)
		for _, domain := range domains {
			if strings.HasSuffix(host, domain) {
				return stateOnFound
			}
		}
		return filter.StateNeutral
	}
}
