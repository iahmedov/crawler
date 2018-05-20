package strategy

import (
	"fmt"

	"github.com/iahmedov/crawler/filter"
	"github.com/iahmedov/crawler/validation"
)

func ValidateConfig(config filter.Config, v *validation.Validator) {
	validateStrategyConfig := func(name string, config filter.StrategyConfig) {
		strategyIfc, ok := config["name"]
		if !ok {
			v.Add(fmt.Errorf("missing name field in %s strategy", name))
			return
		}
		strategyName, ok := strategyIfc.(string)
		if !ok {
			v.Add(fmt.Errorf("strategy name field must be string: %s", name))
			return
		}
		if _, ok := strategies[strategyName]; !ok {
			v.Add(fmt.Errorf("no such %s filter strategy: %s", name, strategyName))
			return
		}
	}
	validateStrategyConfig("link", config.LinkStrategy)
	validateStrategyConfig("content", config.LinkStrategy)
}
