package storage

import (
	"errors"
	"fmt"

	"github.com/iahmedov/crawler/validation"
)

type KV map[string]interface{}
type Config struct {
	Link, Page []KV
}

func ValidateConfigs(config Config, v *validation.Validator) {
	storageConfigValidator := func(configType string, configs []KV, hasFactory func(name string) bool) {
		for _, c := range configs {
			name, ok := c["name"]
			if !ok {
				v.Add(errors.New("missing name field"))
				return
			}

			nameStr, ok := name.(string)
			if !ok {
				v.Add(errors.New("name field should be string value"))
				return
			}

			if !hasFactory(nameStr) {
				v.Add(fmt.Errorf("no %s store with name: %s", configType, nameStr))
				return
			}
		}
	}

	storageConfigValidator("link", config.Link, func(name string) bool {
		_, ok := linkStores[name]
		return ok
	})
	storageConfigValidator("page", config.Page, func(name string) bool {
		_, ok := pageStores[name]
		return ok
	})
}
