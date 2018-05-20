package filter

type FilterConfig map[string]interface{}
type StrategyConfig map[string]interface{}

type Config struct {
	Link            []FilterConfig
	Content         []FilterConfig
	LinkStrategy    StrategyConfig `yaml:"link.strategy"`
	ContentStrategy StrategyConfig `yaml:"content.strategy"`
}
