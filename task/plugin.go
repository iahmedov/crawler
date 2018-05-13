package task

import "fmt"

type StateTransitionerFactory func(Config) StateTransitioner

var transitioners map[string]StateTransitionerFactory

func RegisterStateTransitionerFactory(name string, factory StateTransitionerFactory) {
	if _, ok := transitioners[name]; ok {
		panic(fmt.Sprintf("StateTransitioner already exists with name: %s", name))
	}
	transitioners[name] = factory
}

func BuildStateTransitioner(name string, config Config) StateTransitioner {
	factory, ok := transitioners[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadStateTransitioner(config Config) StateTransitioner {
	if name, ok := config["name"]; ok {
		nameStr, ok := name.(string)
		if ok {
			return BuildStateTransitioner(nameStr, config)
		}
	}

	return nil
}
