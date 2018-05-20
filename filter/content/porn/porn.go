package porn

import (
	"strings"

	"github.com/iahmedov/crawler/filter"
	"github.com/iahmedov/crawler/filter/content"
	"github.com/iahmedov/crawler/model"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	content.RegisterFilterFactory("porn", New)
}

type pornConfig struct {
	Keywords []string `mapstructure:"keywords"`
}

func New(config filter.FilterConfig) (content.Filter, error) {
	fc := &pornConfig{}
	if err := mapstructure.Decode(config, fc); err != nil {
		return nil, errors.Wrap(err, "could not parse porn config")
	}

	return PornFilter(fc.Keywords), nil
}

func PornFilter(keywords []string) content.Filter {
	return func(entry model.CrawlEntry) filter.State {
		txt := entry.Document.Text()
		for _, k := range keywords {
			if strings.Contains(txt, k) {
				return filter.StateHighNegative
			}
		}
		return filter.StateNeutral
	}
}
