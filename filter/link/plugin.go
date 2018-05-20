package link

import (
	"fmt"

	"github.com/iahmedov/crawler/filter"
)

type FilterFactory func(filter.FilterConfig) (Filter, error)

var filters map[string]FilterFactory

func init() {
	filters = map[string]FilterFactory{}
}

func RegisterFilterFactory(name string, factory FilterFactory) {
	if _, ok := filters[name]; ok {
		panic(fmt.Sprintf("filter already exists with name: %s", name))
	}
	filters[name] = factory
}

func BuildFilter(name string, config filter.FilterConfig) (Filter, error) {
	factory, ok := filters[name]
	if !ok {
		return nil, fmt.Errorf("%s link filter not found", name)
	}

	return factory(config)
}

func LoadFilters(configs []filter.FilterConfig) ([]Filter, error) {
	items := []Filter{}
	for _, config := range configs {
		if name, ok := config["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				f, err := BuildFilter(nameStr, config)
				if err != nil {
					return nil, err
				}
				items = append(items, f)
			}
		}
	}

	return items, nil
}
