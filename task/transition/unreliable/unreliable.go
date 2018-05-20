package unreliable

import (
	"github.com/iahmedov/crawler/task"
)

func init() {
	task.RegisterStateTransitionerFactory("unreliable", func(task.Config) (task.StateTransitioner, error) {
		return &unreliable{}, nil
	})
}

type unreliable struct {
}

func (u *unreliable) Transition(task task.Task, state task.TaskState, comment string) error {
	return nil
}
