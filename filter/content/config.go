package content

import (
	"errors"
	"fmt"

	"github.com/iahmedov/crawler/filter"
	"github.com/iahmedov/crawler/validation"
)

func ValidateConfig(config filter.Config, v *validation.Validator) {
	for _, f := range config.Content {
		name, ok := f["name"]
		if !ok {
			v.Add(errors.New("missing name content filter"))
			return
		}

		nameStr, ok := name.(string)
		if !ok {
			v.Add(errors.New("name field should be string value"))
			return
		}

		if _, ok = filters[nameStr]; !ok {
			v.Add(fmt.Errorf("no content filter with name: %s", nameStr))
			return
		}
	}
}
