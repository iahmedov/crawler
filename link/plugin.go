package link

import (
	"fmt"
)

type ExtractorFactory func(Config) Extractor

var extractors map[string]ExtractorFactory

func RegisterExtractorFactory(name string, factory ExtractorFactory) {
	if _, ok := extractors[name]; ok {
		panic(fmt.Sprintf("link extractor already exists with name: %s", name))
	}
	extractors[name] = factory
}

func BuildExtractor(name string, config Config) Extractor {
	factory, ok := extractors[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadExtractors(configs []Config) []Extractor {
	items := []Extractor{}
	for _, config := range configs {
		if name, ok := config["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				items = append(items, BuildExtractor(nameStr, config))
			}
		}
	}

	return items
}
