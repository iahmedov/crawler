package storage

import (
	"io"

	"github.com/iahmedov/crawler/model"
)

type LinkStore interface {
	io.Closer

	SaveLinks(links []model.Link, parent model.Link) error
	HasLink(link model.Link) (bool, error)
	Link(model.URL) (*model.Link, error)
}

type PageStore interface {
	io.Closer

	Page(model.URL) (*model.Page, error)
	PageWithEntries(model.URL) (*model.Page, error)
	AppendCrawlEntry(*model.CrawlEntry) error
}
