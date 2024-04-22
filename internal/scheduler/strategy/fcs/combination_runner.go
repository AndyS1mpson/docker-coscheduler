package fcs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/combination"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/pointer"
)

// CombinationRunner описывает логику взаимодействия с нодами для эмпирической оценки задач
type CombinationRunner struct {
	node                models.Node
	client              nodeClient
	taskCombintationNum int64
	measurementTime     time.Duration
}

// NewCombinationRunner конструктор для CombinationRunner
func NewCombinationRunner(node models.Node, client nodeClient, taskCombinationNum int64, measurementTime time.Duration) *CombinationRunner {
	return &CombinationRunner{
		node:                node,
		client:              client,
		taskCombintationNum: taskCombinationNum,
		measurementTime:     measurementTime,
	}
}

// Run выполняет 1 комбинацию на ноде и замеряет время выполнения
func (nw *CombinationRunner) Run(ctx context.Context, taskCombo combination.Combination) (*time.Duration, error) {
	var total time.Duration

	taskTimes := make(chan (time.Duration), len(taskCombo.Tasks))
	errChan := make(chan error, 1)

	g, ctx := errgroup.WithContext(ctx)

	var mu sync.Mutex

	for task := range taskCombo.Tasks {
		if task.Node.Host != nw.node.Host || task.Node.Port != nw.node.Port {
			continue
		}

		g.Go(func() error {
			var (
				err      error
				duration time.Duration
			)

			oldTask := task

			duration, err = nw.executeTask(ctx, &task)
			if err != nil {
				return err
			}

			taskTimes <- duration

			mu.Lock()
			defer mu.Unlock()

			delete(taskCombo.Tasks, oldTask)
			taskCombo.Tasks[task] = struct{}{}

			return nil
		})
	}

	go func() {
		defer close(taskTimes)
		if err := g.Wait(); err != nil {
			log.Error(err, log.Data{})

			errChan <- err

			return
		}
	}()

	for duration := range taskTimes {
		total += duration
	}

	select {
	case err := <-errChan:
		return nil, err
	default:
		return &total, nil
	}
}

func (nw *CombinationRunner) executeTask(ctx context.Context, task *models.Task) (time.Duration, error) {
	containerID, err := nw.client.CreateTask(ctx, task.ImageID, pointer.Get(task.Config).CPUs)
	if err != nil {
		return 0, fmt.Errorf("create task with id %s: %w", task.ID, err)
	}

	if task.Config == nil {
		task.Config = &models.Config{}
	}

	task.Config.ContainerID = containerID

	err = nw.client.StartTask(ctx, containerID)
	if err != nil {
		return 0, fmt.Errorf("start task with id %s: %w", task.ID, err)
	}

	err = nw.client.PauseTask(ctx, containerID)
	if err != nil {
		return 0, fmt.Errorf("pause task with id %s: %w", task.ID, err)
	}

	duration, err := nw.client.MeasureTaskSpeed(ctx, containerID, pointer.Get(pointer.Get(task.Config).CPUs), nw.measurementTime)
	if err != nil {
		return 0, fmt.Errorf("measure task speed with id %s: %w", task.ID, err)
	}

	return duration, nil
}
