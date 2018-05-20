package model

import (
	"bytes"
	"io"
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

func (d Document) Reader() io.Reader {
	return bytes.NewReader(d)
}

func NewCrawlEntry(parent url.URL, depth uint32) *CrawlEntry {
	return &CrawlEntry{
		Parent:    parent,
		Error:     nil,
		Depth:     depth,
		CreatedAt: time.Now().UTC(),
	}
}
