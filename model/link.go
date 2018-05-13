package model

import (
	"net/url"

	"github.com/iahmedov/crawler/task"
)

type URL = url.URL

// let's not do it in this version
// type URLCleaner func(u URL) URL

type Link struct {
	URI, Parent   URL // both links must be absolute URI path
	RelativeDepth uint32
}

func (nk *Link) URL() url.URL {
	return nk.URI
}

func (nk *Link) Depth() uint32 {
	return nk.RelativeDepth
}

var _ task.Task = (*Link)(nil)
