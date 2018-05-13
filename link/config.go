package link

import (
	"errors"

	"github.com/iahmedov/crawler/validation"
)

type Config map[string]interface{}

func ValidateConfigs(configs []Config, v *validation.Validator) {
	for _, c := range configs {
		if _, ok := c["name"]; !ok {
			v.Add(errors.New("missing name field in link extractors"))
			return
		}
	}
}
