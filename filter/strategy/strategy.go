package strategy

import (
	"math"

	"github.com/iahmedov/crawler/filter"
)

func init() {
	RegisterStrategyBuilderFactory("cumulative", func(filter.StrategyConfig) (StrategyBuilder, error) {
		return strategyBuilderFunc(func() Strategy {
			return &cumulative{}
		}), nil
	})

	RegisterStrategyBuilderFactory("negative_high", func(filter.StrategyConfig) (StrategyBuilder, error) {
		return strategyBuilderFunc(func() Strategy {
			return &negativeHigh{
				negative: filter.StatePositive,
			}
		}), nil
	})
}

type Strategy interface {
	AddFilterResult(filter.State)
	Decision() filter.State
}

type StrategyBuilder interface {
	Build() Strategy
}

type strategyBuilderFunc func() Strategy

func (s strategyBuilderFunc) Build() Strategy {
	return s()
}

type cumulative struct {
	results []filter.State
}

func (c *cumulative) AddFilterResult(result filter.State) {
	c.results = append(c.results, result)
}

func (c *cumulative) Decision() filter.State {
	sum := 0
	for _, v := range c.results {
		sum += int(v)
	}

	mean := sum / len(c.results)
	return filter.State(mean)
}

type negativeHigh struct {
	negative filter.State
}

func maxValueOf(v1, v2 filter.State) filter.State {
	if math.Abs(float64(v1)) < math.Abs(float64(v2)) {
		return v2
	}
	return v1
}

func (n *negativeHigh) AddFilterResult(result filter.State) {
	if result.IsNegative() {
		n.negative = maxValueOf(n.negative, result)
	}
}

func (n *negativeHigh) Decision() filter.State {
	return n.negative
}
