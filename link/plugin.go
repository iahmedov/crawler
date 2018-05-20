package link

import (
	"fmt"
)

type ExtractorFactory func(Config) (Extractor, error)

var extractors map[string]ExtractorFactory

func init() {
	extractors = map[string]ExtractorFactory{}
}

func RegisterExtractorFactory(name string, factory ExtractorFactory) {
	if _, ok := extractors[name]; ok {
		panic(fmt.Sprintf("link extractor already exists with name: %s", name))
	}
	extractors[name] = factory
}

func BuildExtractor(name string, config Config) (Extractor, error) {
	factory, ok := extractors[name]
	if !ok {
		return nil, fmt.Errorf("%s extractor not found", name)
	}

	return factory(config)
}

func LoadExtractors(configs []Config) ([]Extractor, error) {
	items := []Extractor{}
	for _, config := range configs {
		if name, ok := config["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				extractor, err := BuildExtractor(nameStr, config)
				if err != nil {
					return nil, err
				}
				items = append(items, extractor)
			}
		}
	}

	return items, nil
}
