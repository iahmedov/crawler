package model

import (
	"net/http"
	"net/url"
	"time"
)

type Project uint32

type Page struct {
	CrawlEntries []*CrawlEntry
	Link         URL
	CreatedAt    time.Time
	PageStatus   PageStatus
	Project      Project // not used yet
}

type CrawlEntry struct {
	Parent     url.URL
	StatusCode int
	Header     http.Header
	Error      error
	Document   Document
	Links      []URL
	Depth      uint32
	CreatedAt  time.Time
}

type Document []byte
type PageStatus int

const (
	PageStatusUnknown PageStatus = iota
	PageStatusOK
	PageStatusBlocked
	PageStatusStale
)

func (d Document) Text() string {
	return string(d)
}

func NewCrawlEntry(parent url.URL) *CrawlEntry {
	return &CrawlEntry{
		Parent:    parent,
		Error:     nil,
		CreatedAt: time.Now().UTC(),
	}
}
