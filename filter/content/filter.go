package content

import (
	"github.com/iahmedov/crawler/filter"
	"github.com/iahmedov/crawler/model"
)

type Filter func(model.CrawlEntry) filter.State
