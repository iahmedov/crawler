package queue

import (
	"context"

	"github.com/iahmedov/crawler/task"
)

type Queue interface {
	Tasks(ctx context.Context) <-chan task.Task
	Put([]task.Task) error

	SetStateTransitioner(task.StateTransitioner)
}
