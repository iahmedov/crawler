package strategy

import (
	"fmt"

	"github.com/iahmedov/crawler/filter"
)

type StrategyFactory func(filter.StrategyConfig) Strategy

var strategies map[string]StrategyFactory

func RegisterStrategyFactory(name string, factory StrategyFactory) {
	if _, ok := strategies[name]; ok {
		panic(fmt.Sprintf("strategy already exists with name: %s", name))
	}
	strategies[name] = factory
}

func BuildStrategy(name string, config filter.StrategyConfig) Strategy {
	factory, ok := strategies[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadStrategy(config filter.StrategyConfig) Strategy {
	if name, ok := config["name"]; ok {
		nameStr, ok := name.(string)
		if ok {
			return BuildStrategy(nameStr, config)
		}
	}

	return nil
}
