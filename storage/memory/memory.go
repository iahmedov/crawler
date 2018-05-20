package memory

import (
	"sync"

	"github.com/iahmedov/crawler/model"
	"github.com/iahmedov/crawler/storage"
)

func init() {
	storage.RegisterLinkStoreFactory("memory", NewLinkStore)
	storage.RegisterPageStoreFactory("memory", NewPageStore)
}

var singletonStore *memoryStore
var mtxStore sync.Mutex

type memoryStore struct {
	mtxUrls  sync.Mutex
	mtxPages sync.Mutex
	urls     map[string]model.Link
	pages    map[string]*model.Page
}

func NewLinkStore(config storage.KV) (storage.LinkStore, error) {
	mtxStore.Lock()
	defer mtxStore.Unlock()
	if singletonStore != nil {
		return singletonStore, nil
	}

	singletonStore := &memoryStore{
		urls:  map[string]model.Link{},
		pages: map[string]*model.Page{},
	}
	return singletonStore, nil
}

func NewPageStore(config storage.KV) (storage.PageStore, error) {
	return nil, nil
}

func (m *memoryStore) Close() error {
	return nil
}

func (m *memoryStore) SaveLinks(links []model.Link, parent model.Link) error {
	m.mtxUrls.Lock()
	defer m.mtxUrls.Unlock()
	depth := parent.Depth() + 1

	for _, nk := range links {
		u := nk.URL()
		urlString := u.String()
		stored, ok := m.urls[urlString]
		if !ok {
			nk.RelativeDepth = depth
			m.urls[urlString] = nk
			continue
		}

		if depth < stored.RelativeDepth {
			stored.RelativeDepth = depth
			stored.Parent = parent.URL()
		}
	}

	return nil
}

func (m *memoryStore) HasLink(link model.Link) (bool, error) {
	item, _ := m.Link(link.URL())
	return item != nil, nil
}

func (m *memoryStore) Link(u model.URL) (*model.Link, error) {
	m.mtxUrls.Lock()
	defer m.mtxUrls.Unlock()
	link, _ := m.urls[u.String()]
	return &link, nil
}

func (m *memoryStore) Page(u model.URL) (*model.Page, error) {
	m.mtxPages.Lock()
	defer m.mtxPages.Unlock()
	p, _ := m.pages[u.String()]
	return p, nil
}

func (m *memoryStore) PageWithEntries(u model.URL) (*model.Page, error) {
	return m.Page(u)
}

func (m *memoryStore) AppendCrawlEntry(entry *model.CrawlEntry) error {
	p, err := m.Page(entry.Parent)
	if err != nil {
		return err
	}

	p.CrawlEntries = append(p.CrawlEntries, entry)
	return nil
}
