package strategy

import (
	"github.com/iahmedov/crawler/filter"
)

type Strategy interface {
	AddFilterResult(filter.State)
	Decision() filter.State
}
