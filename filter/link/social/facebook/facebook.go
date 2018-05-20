package facebook

import (
	"net/url"
	"strings"

	"github.com/iahmedov/crawler/filter"
	filterlink "github.com/iahmedov/crawler/filter/link"
)

func init() {
	filterlink.RegisterFilterFactory("facebook", New)
}

func New(config filter.FilterConfig) (filterlink.Filter, error) {
	return Facebook(), nil
}

func Facebook() filterlink.Filter {
	return func(u url.URL, depth uint32) filter.State {
		host := u.Hostname()
		host = strings.ToLower(host)
		// afb.com?
		if strings.Contains(host, "facebook.com") || strings.HasSuffix(host, "fb.com") {
			return filter.StateHighNegative
		}
		return filter.StateNeutral
	}
}
