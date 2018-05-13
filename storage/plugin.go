package storage

import (
	"fmt"
)

type LinkStoreFactory func(KV) LinkStore
type PageStoreFactory func(KV) PageStore

var linkStores map[string]LinkStoreFactory
var pageStores map[string]PageStoreFactory

func RegisterLinkStoreFactory(name string, factory LinkStoreFactory) {
	if _, ok := linkStores[name]; ok {
		panic(fmt.Sprintf("link store already exists with name: %s", name))
	}
	linkStores[name] = factory
}

func BuildLinkStore(name string, config KV) LinkStore {
	factory, ok := linkStores[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadLinkStores(config Config) []LinkStore {
	items := []LinkStore{}
	for _, c := range config.Link {
		if name, ok := c["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				items = append(items, BuildLinkStore(nameStr, c))
			}
		}
	}

	return items
}

func RegisterPageStoreFactory(name string, factory PageStoreFactory) {
	if _, ok := pageStores[name]; ok {
		panic(fmt.Sprintf("page store already exists with name: %s", name))
	}
	pageStores[name] = factory
}

func BuildPageStore(name string, config KV) PageStore {
	factory, ok := pageStores[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadPageStores(config Config) []PageStore {
	items := []PageStore{}
	for _, c := range config.Page {
		if name, ok := c["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				items = append(items, BuildPageStore(nameStr, c))
			}
		}
	}

	return items
}
