package html

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	link "github.com/iahmedov/crawler/link"
	"github.com/iahmedov/crawler/model"
)

func init() {
	link.RegisterExtractorFactory("html", NewHtmlLinkExtractor)
}

type HtmlLinkExtractor struct {
}

func NewHtmlLinkExtractor(config link.Config) (link.Extractor, error) {
	return &HtmlLinkExtractor{}, nil
}

func (h *HtmlLinkExtractor) Extract(entry model.CrawlEntry) []url.URL {
	doc, err := goquery.NewDocumentFromReader(entry.Document.Reader())
	if err != nil {
		return []url.URL{}
	}

	urls := []url.URL{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			u, err := url.Parse(href)
			if err != nil {
				urls = append(urls, *u)
			}
		}
	})

	return urls
}
