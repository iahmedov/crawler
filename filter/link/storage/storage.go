package storage

import (
	"net/url"

	"github.com/iahmedov/crawler/filter"
	filterlink "github.com/iahmedov/crawler/filter/link"
	"github.com/iahmedov/crawler/model"
	"github.com/iahmedov/crawler/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	filterlink.RegisterFilterFactory("unique_from_storage", New)
}

type storageConfig struct {
	Storage storage.KV `mapstructure:"storage"`
}

func New(config filter.FilterConfig) (filterlink.Filter, error) {
	var sc storageConfig
	if err := mapstructure.Decode(config, &sc); err != nil {
		return nil, err
	}
	nameIfc, ok := sc.Storage["name"]
	if !ok {
		return nil, errors.New("storage doesnt contain name field")
	}
	name, ok := nameIfc.(string)
	if !ok {
		return nil, errors.New("name should be string field")
	}

	linkStore, err := storage.BuildLinkStore(name, sc.Storage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build link store")
	}

	return FromStorage(linkStore), nil
}

func FromStorage(linkStore storage.LinkStore) filterlink.Filter {
	return func(u url.URL, depth uint32) filter.State {
		has, err := linkStore.HasLink(model.Link{
			URI: u,
		})
		if err != nil {
			return filter.StateNeutral
		}
		if has {
			return filter.StateNegative
		}
		return filter.StatePositive
	}
}
