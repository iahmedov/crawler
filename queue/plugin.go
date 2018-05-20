package queue

import (
	"fmt"

	"github.com/iahmedov/crawler/validation"
)

type QueueFactory func(Config) (Queue, error)

var queues map[string]QueueFactory

func init() {
	queues = map[string]QueueFactory{}
}

func RegisterQueueFactory(name string, factory QueueFactory) {
	if _, ok := queues[name]; ok {
		panic(fmt.Sprintf("queue already exists with name: %s", name))
	}
	queues[name] = factory
}

func BuildQueue(name string, config Config) (Queue, error) {
	factory, ok := queues[name]
	if !ok {
		return nil, fmt.Errorf("%s queue not found", name)
	}

	return factory(config)
}

func LoadQueue(config Config) (Queue, error) {
	v := validation.NewValidator("queue")
	ValidateConfig(config, v)
	if v.HasError() {
		return nil, v
	}

	name, _ := config["name"]
	nameStr, _ := name.(string)
	return BuildQueue(nameStr, config)
}
