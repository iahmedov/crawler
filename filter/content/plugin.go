package content

import (
	"fmt"

	"github.com/iahmedov/crawler/filter"
)

type FilterFactory func(filter.FilterConfig) Filter

var filters map[string]FilterFactory

func RegisterFilterFactory(name string, factory FilterFactory) {
	if _, ok := filters[name]; ok {
		panic(fmt.Sprintf("filter already exists with name: %s", name))
	}
	filters[name] = factory
}

func BuildFilter(name string, config filter.FilterConfig) Filter {
	factory, ok := filters[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadFilters(configs []filter.FilterConfig) []Filter {
	items := []Filter{}
	for _, config := range configs {
		if name, ok := config["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				items = append(items, BuildFilter(nameStr, config))
			}
		}
	}

	return items
}
