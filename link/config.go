package link

import (
	"errors"
	"fmt"

	"github.com/iahmedov/crawler/validation"
)

type Config map[string]interface{}

func ValidateConfigs(configs []Config, v *validation.Validator) {
	for _, c := range configs {
		name, ok := c["name"]
		if !ok {
			v.Add(errors.New("missing name field in link extractors"))
			return
		}

		nameStr, ok := name.(string)
		if !ok {
			v.Add(errors.New("name field should be string value"))
			return
		}

		if _, ok = extractors[nameStr]; !ok {
			v.Add(fmt.Errorf("no extractor with name: %s", nameStr))
			return
		}
	}
}
