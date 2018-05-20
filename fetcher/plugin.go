package fetcher

import (
	"fmt"
)

type MiddlewareFactory func(Config) (Middleware, error)

var middlewares map[string]MiddlewareFactory

func init() {
	middlewares = map[string]MiddlewareFactory{}
}

func RegisterMiddlewareFactory(name string, factory MiddlewareFactory) {
	if _, ok := middlewares[name]; ok {
		panic(fmt.Sprintf("middleware already exists with name: %s", name))
	}
	middlewares[name] = factory
}

func BuildMiddleware(name string, config Config) (Middleware, error) {
	factory, ok := middlewares[name]
	if !ok {
		return nil, fmt.Errorf("%s middleware not found", name)
	}

	return factory(config)
}

func LoadMiddlewares(configs []Config) ([]Middleware, error) {
	items := []Middleware{}
	for _, config := range configs {
		if name, ok := config["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				middle, err := BuildMiddleware(nameStr, config)
				if err != nil {
					return nil, err
				}
				items = append(items, middle)
			}
		}
	}

	return items, nil
}
