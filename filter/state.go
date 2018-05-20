package filter

type State int

const (
	StateHighNegative State = -2
	StateNegative     State = -1
	StateNeutral      State = 0
	StatePositive     State = 1
	StateHighPositive State = 2
)

func (s State) IsPositive() bool {
	return s > 0
}

func (s State) IsNegative() bool {
	return s < 0
}

func (s State) IsNeutral() bool {
	return s == 0
}
