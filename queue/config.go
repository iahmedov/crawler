package queue

import (
	"errors"

	"github.com/iahmedov/crawler/validation"
)

type Config map[string]interface{}

func ValidateConfig(config Config, v *validation.Validator) {
	if _, ok := config["name"]; !ok {
		v.Add(errors.New("missing name field"))
		return
	}
}
