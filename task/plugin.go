package task

import (
	"fmt"

	"github.com/iahmedov/crawler/validation"
)

type StateTransitionerFactory func(Config) (StateTransitioner, error)

var transitioners map[string]StateTransitionerFactory

func init() {
	transitioners = map[string]StateTransitionerFactory{}
}

func RegisterStateTransitionerFactory(name string, factory StateTransitionerFactory) {
	if _, ok := transitioners[name]; ok {
		panic(fmt.Sprintf("StateTransitioner already exists with name: %s", name))
	}
	transitioners[name] = factory
}

func BuildStateTransitioner(name string, config Config) (StateTransitioner, error) {
	factory, ok := transitioners[name]
	if !ok {
		return nil, fmt.Errorf("%s transitioner not found", name)
	}

	return factory(config)
}

func LoadStateTransitioner(config Config) (StateTransitioner, error) {
	v := validation.NewValidator("task")
	ValidateConfig(config, v)
	if v.HasError() {
		return nil, nil
	}

	name, _ := config["name"]
	nameStr, _ := name.(string)
	return BuildStateTransitioner(nameStr, config)
}
