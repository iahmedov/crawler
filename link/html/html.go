package html

import (
	"net/url"

	link "github.com/iahmedov/crawler/link"
	"github.com/iahmedov/crawler/model"
)

func Init() {
	link.RegisterExtractorFactory("html", NewHtmlLinkExtractor)
}

type HtmlLinkExtractor struct {
}

func NewHtmlLinkExtractor(config link.Config) link.LinkExtractor {
	return &HtmlLinkExtractor{}
}

func (h *HtmlLinkExtractor) Extract(entry model.CrawlEntry) []url.URL {
	return []url.URL{}
}

var _ link.LinkExtractor = (*HtmlLinkExtractor)(nil)
