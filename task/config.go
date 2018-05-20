package task

import (
	"errors"
	"fmt"

	"github.com/iahmedov/crawler/validation"
)

type Config map[string]interface{}

func ValidateConfig(config Config, v *validation.Validator) {
	name, ok := config["name"]
	if !ok {
		v.Add(errors.New("missing name field in task.state.transitioner"))
		return
	}

	nameStr, ok := name.(string)
	if !ok {
		v.Add(errors.New("name field should be string value"))
		return
	}
	_, ok = transitioners[nameStr]
	if !ok {
		v.Add(fmt.Errorf("no state transitioner with name: %s", nameStr))
		return
	}
}
