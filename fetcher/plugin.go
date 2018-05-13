package fetcher

import (
	"fmt"
)

type MiddlewareFactory func(Config) Middleware

var middlewares map[string]MiddlewareFactory

func RegisterMiddlewareFactory(name string, factory MiddlewareFactory) {
	if _, ok := middlewares[name]; ok {
		panic(fmt.Sprintf("middleware already exists with name: %s", name))
	}
	middlewares[name] = factory
}

func BuildMiddleware(name string, config Config) Middleware {
	factory, ok := middlewares[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadMiddlewares(configs []Config) []Middleware {
	items := []Middleware{}
	for _, config := range configs {
		if name, ok := config["name"]; ok {
			nameStr, ok := name.(string)
			if ok {
				items = append(items, BuildMiddleware(nameStr, config))
			}
		}
	}

	return items
}
