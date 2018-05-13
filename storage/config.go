package storage

import (
	"errors"

	"github.com/iahmedov/crawler/validation"
)

type KV map[string]interface{}
type Config struct {
	Link, Page []KV
}

func ValidateConfigs(config Config, v *validation.Validator) {
	configs := append(config.Link, config.Page...)
	for _, c := range configs {
		if _, ok := c["name"]; !ok {
			v.Add(errors.New("missing name field"))
			return
		}
	}
}
