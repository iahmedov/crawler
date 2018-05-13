package queue

import (
	"fmt"
)

type QueueFactory func(Config) Queue

var queues map[string]QueueFactory

func RegisterQueueFactory(name string, factory QueueFactory) {
	if _, ok := queues[name]; ok {
		panic(fmt.Sprintf("queue already exists with name: %s", name))
	}
	queues[name] = factory
}

func BuildQueue(name string, config Config) Queue {
	factory, ok := queues[name]
	if !ok {
		return nil
	}

	return factory(config)
}

func LoadQueue(config Config) Queue {
	if name, ok := config["name"]; ok {
		nameStr, ok := name.(string)
		if ok {
			return BuildQueue(nameStr, config)
		}
	}

	return nil
}
