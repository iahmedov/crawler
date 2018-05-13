package unreliable

import "github.com/iahmedov/crawler/task"

func Init() {
	task.RegisterStateTransitionerFactory("unreliable", func(task.Config) task.StateTransitioner {
		return &unreliable{}
	})
}

type unreliable struct {
}

func (u *unreliable) Transition(task task.Task, state task.TaskState) error {
	return nil
}
