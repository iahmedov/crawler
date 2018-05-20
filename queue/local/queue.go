package local

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/iahmedov/crawler/queue"
	"github.com/iahmedov/crawler/task"
	"github.com/mitchellh/mapstructure"
)

type localQueueConfig struct {
	BufferSize   int `mapstructure:"buffer_size"`
	LinkWaitTime int `mapstructure:"link_wait_time"`
}

type localQueue struct {
	taskTransitioner task.StateTransitioner
	buffer           []task.Task
	mtxBuffer        sync.Mutex
	publisher        chan task.Task
	config           localQueueConfig
}

func init() {
	queue.RegisterQueueFactory("local", New)
}

func New(config queue.Config) (queue.Queue, error) {
	q := &localQueue{}
	if err := mapstructure.Decode(config, &q.config); err != nil {
		return nil, errors.Wrap(err, "could not decode config")
	}
	return q, nil
}

func (q *localQueue) Run(ctx context.Context) error {
	go func() {
		ctxDone := ctx.Done()
		var nextTask task.Task
	loop:
		for {
			select {
			case <-ctxDone:
				break loop
			default:
			}

			q.mtxBuffer.Lock()
			if len(q.buffer) > 0 {
				nextTask = q.buffer[0]
				q.buffer = q.buffer[1:] // danger
			} else {
				nextTask = nil
			}
			q.mtxBuffer.Unlock()

			if nextTask == nil {
				time.Sleep(time.Duration(q.config.LinkWaitTime) * time.Millisecond)
				continue
			}

			q.publish(nextTask)
		}
	}()
	return nil
}

func (q *localQueue) publish(next task.Task) {
	q.publisher <- next
}

func (q *localQueue) Tasks(ctx context.Context) <-chan task.Task {
	taskChan := make(chan task.Task, 1)
	go func() {
		ctxDone := ctx.Done()
	loop:
		for {
			select {
			case t := <-q.publisher:
				taskChan <- t
			case <-ctxDone:
				break loop
			}
		}
		close(taskChan)
	}()
	return taskChan
}

func (q *localQueue) Put(tasks []task.Task) error {
	q.mtxBuffer.Lock()
	defer q.mtxBuffer.Unlock()

	if len(q.buffer) > q.config.BufferSize {
		return errors.New("buffer is full: all tasks dropped")
	}

	if len(q.buffer)+len(tasks) < q.config.BufferSize {
		return errors.New("buffer is almost full: no space for all tasks")
	}

	for _, t := range tasks {
		q.buffer = append(q.buffer, t)
		q.taskTransitioner.Transition(t, task.TaskStateInitial, "initial")
	}

	return nil
}

func (q *localQueue) SetStateTransitioner(transitioner task.StateTransitioner) {
	q.taskTransitioner = transitioner
}
