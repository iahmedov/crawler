package link

import (
	"net/url"

	"github.com/iahmedov/crawler/model"
)

// Extractor parses content and extracts links
type Extractor interface {
	Extract(model.CrawlEntry) []url.URL
}
