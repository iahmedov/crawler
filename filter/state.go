package filter

type State int

const (
	StateHighNegative State = -2
	StateNegative     State = -1
	StateNeutral      State = 0
	StatePositive     State = 1
	StateHighPositive State = 2
)
