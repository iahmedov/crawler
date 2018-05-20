package strategy

import (
	"fmt"

	"github.com/iahmedov/crawler/filter"
)

type StrategyBuilderFactory func(filter.StrategyConfig) (StrategyBuilder, error)

var strategies map[string]StrategyBuilderFactory

func init() {
	strategies = map[string]StrategyBuilderFactory{}
}

func RegisterStrategyBuilderFactory(name string, factory StrategyBuilderFactory) {
	if _, ok := strategies[name]; ok {
		panic(fmt.Sprintf("strategy already exists with name: %s", name))
	}
	strategies[name] = factory
}

func BuildStrategyBuilder(name string, config filter.StrategyConfig) (StrategyBuilder, error) {
	factory, ok := strategies[name]
	if !ok {
		return nil, fmt.Errorf("%s strategy not found", name)
	}

	return factory(config)
}

func LoadStrategyBuilder(config filter.StrategyConfig) (StrategyBuilder, error) {
	if name, ok := config["name"]; ok {
		nameStr, ok := name.(string)
		if ok {
			return BuildStrategyBuilder(nameStr, config)
		}
	}

	return nil, nil
}
