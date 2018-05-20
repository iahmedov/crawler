package queue

import (
	"errors"
	"fmt"

	"github.com/iahmedov/crawler/validation"
)

type Config map[string]interface{}

func ValidateConfig(config Config, v *validation.Validator) {
	name, ok := config["name"]
	if !ok {
		v.Add(errors.New("missing name field"))
		return
	}

	nameStr, ok := name.(string)
	if !ok {
		v.Add(errors.New("name field should be string value"))
		return
	}

	if _, ok = queues[nameStr]; !ok {
		v.Add(fmt.Errorf("no queue with name: %s", nameStr))
		return
	}
}
