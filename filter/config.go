package filter

import (
	"errors"

	"github.com/iahmedov/crawler/validation"
)

type FilterConfig map[string]interface{}
type StrategyConfig map[string]interface{}

type Config struct {
	Link            []FilterConfig
	Content         []FilterConfig
	LinkStrategy    StrategyConfig `yaml:"link.strategy"`
	ContentStrategy StrategyConfig `yaml:"content.strategy"`
}

func ValidateConfig(config Config, v *validation.Validator) {
	filters := append(config.Link, config.Content...)
	for _, f := range filters {
		if _, ok := f["name"]; !ok {
			v.Add(errors.New("missing name field in link/content filters"))
			return
		}
	}

	if _, ok := config.LinkStrategy["name"]; !ok {
		v.Add(errors.New("missing name field in link.strategy"))
		return
	}

	if _, ok := config.ContentStrategy["name"]; !ok {
		v.Add(errors.New("missing name field in content.strategy"))
		return
	}
}
