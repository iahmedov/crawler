package queue

import (
	"context"

	"github.com/iahmedov/crawler/task"
)

type Queue interface {
	Run(ctx context.Context) error
	Tasks(ctx context.Context) <-chan task.Task
	Put([]task.Task) error

	SetStateTransitioner(task.StateTransitioner)
}
