package task

import "net/url"

type StateTransitioner interface {
	Transition(task Task, state TaskState, comment string) error
}

type Task interface {
	URL() url.URL
	Depth() uint32
}

type TaskState uint32

const (
	TaskStateInitial TaskState = iota
	TaskStateProcessing
	TaskStateProcessed
	TaskStateFailed
	TaskStateDropped
)
