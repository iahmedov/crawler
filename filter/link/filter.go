package link

import (
	"net/url"

	"github.com/iahmedov/crawler/filter"
)

type Filter func(link url.URL, depth uint32) filter.State
