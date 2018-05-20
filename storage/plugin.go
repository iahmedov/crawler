package storage

import (
	"fmt"
)

type LinkStoreFactory func(KV) (LinkStore, error)
type PageStoreFactory func(KV) (PageStore, error)

var linkStores map[string]LinkStoreFactory
var pageStores map[string]PageStoreFactory

func init() {
	linkStores = map[string]LinkStoreFactory{}
	pageStores = map[string]PageStoreFactory{}
}

func RegisterLinkStoreFactory(name string, factory LinkStoreFactory) {
	if _, ok := linkStores[name]; ok {
		panic(fmt.Sprintf("link store already exists with name: %s", name))
	}
	linkStores[name] = factory
}

func BuildLinkStore(name string, config KV) (LinkStore, error) {
	factory, ok := linkStores[name]
	if !ok {
		return nil, fmt.Errorf("%s link store not found", name)
	}

	return factory(config)
}

func LoadLinkStores(config Config) ([]LinkStore, error) {
	items := []LinkStore{}
	for _, c := range config.Link {
		if name, ok := c["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				s, err := BuildLinkStore(nameStr, c)
				if err != nil {
					return nil, err
				}
				items = append(items, s)
			}
		}
	}

	return items, nil
}

func RegisterPageStoreFactory(name string, factory PageStoreFactory) {
	if _, ok := pageStores[name]; ok {
		panic(fmt.Sprintf("page store already exists with name: %s", name))
	}
	pageStores[name] = factory
}

func BuildPageStore(name string, config KV) (PageStore, error) {
	factory, ok := pageStores[name]
	if !ok {
		return nil, fmt.Errorf("%s page store not found", name)
	}

	return factory(config)
}

func LoadPageStores(config Config) ([]PageStore, error) {
	items := []PageStore{}
	for _, c := range config.Page {
		if name, ok := c["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				s, err := BuildPageStore(nameStr, c)
				if err != nil {
					return nil, err
				}
				items = append(items, s)
			}
		}
	}

	return items, nil
}
