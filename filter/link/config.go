package link

import (
	"errors"
	"fmt"

	"github.com/iahmedov/crawler/filter"
	"github.com/iahmedov/crawler/validation"
)

func ValidateConfig(config filter.Config, v *validation.Validator) {
	for _, f := range config.Link {
		name, ok := f["name"]
		if !ok {
			v.Add(errors.New("missing name link filter"))
			return
		}

		nameStr, ok := name.(string)
		if !ok {
			v.Add(errors.New("name field should be string value"))
			return
		}

		if _, ok = filters[nameStr]; !ok {
			v.Add(fmt.Errorf("no link filter with name: %s", nameStr))
			return
		}
	}
}
