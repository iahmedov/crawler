package crawler

import (
	"github.com/iahmedov/crawler/fetcher"
	"github.com/iahmedov/crawler/filter"
	filtercontent "github.com/iahmedov/crawler/filter/content"
	filterlink "github.com/iahmedov/crawler/filter/link"
	filterstrategy "github.com/iahmedov/crawler/filter/strategy"
	"github.com/iahmedov/crawler/link"
	"github.com/iahmedov/crawler/queue"
	"github.com/iahmedov/crawler/storage"
	"github.com/iahmedov/crawler/task"
	"github.com/iahmedov/crawler/validation"
)

type Config struct {
	Workers         int              `yaml:"workers"`
	Queue           queue.Config     `yaml:"queue"`
	StateTransition task.Config      `yaml:"task.state.transitioner"`
	Filters         filter.Config    `yaml:"filters"`
	Fetchers        []fetcher.Config `yaml:"fetchers"`
	Extractors      []link.Config    `yaml:"extractors"`
	Storage         storage.Config   `yaml:"storages"`
}

func ValidateConfig(config Config) *validation.Validator {
	v := validation.NewValidator("config")
	queue.ValidateConfig(config.Queue, v.WithContext("queue"))
	task.ValidateConfig(config.StateTransition, v.WithContext("task"))

	filterlink.ValidateConfig(config.Filters, v.WithContext("filter.link"))
	filtercontent.ValidateConfig(config.Filters, v.WithContext("filter.content"))
	filterstrategy.ValidateConfig(config.Filters, v.WithContext("filter.strategy"))

	fetcher.ValidateConfigs(config.Fetchers, v.WithContext("fetcher"))
	link.ValidateConfigs(config.Extractors, v.WithContext("extractor"))
	storage.ValidateConfigs(config.Storage, v.WithContext("storage"))

	if v.HasError() {
		return v
	}
	return nil
}
