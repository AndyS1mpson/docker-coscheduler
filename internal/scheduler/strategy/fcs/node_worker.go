package fcs

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/combination"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
)

// NodeWorker описывает логику взаимодействия с нодами для эмпирической оценки задач
type NodeWorker struct {
	node                models.Node
	client              nodeClient
	taskCombintationNum int64
}

// NewNodeWorker конструктор для NodeWorker
func NewNodeWorker(node models.Node, client nodeClient, taskCombinationNum int64) *NodeWorker {
	return &NodeWorker{
		node:                node,
		client:              client,
		taskCombintationNum: taskCombinationNum,
	}
}

// preStage выполняет 1 комбинацию на ноде и замеряет время выполнения
func (nw *NodeWorker) preStage(ctx context.Context, taskCombo combination.Combination) (*time.Duration, error) {
	var total time.Duration

	taskTimes := make(chan (time.Duration), len(taskCombo.Tasks))
	errChan := make(chan error, 1)

	g, ctx := errgroup.WithContext(ctx)

	for task := range taskCombo.Tasks {
		if task.Node.Host != nw.node.Host || task.Node.Port != nw.node.Port {
			continue
		}

		g.Go(func() error {
			var err error
			var duration *time.Duration
			duration, err = nw.executeTask(ctx, task)
			if err != nil {
				return err
			}

			taskTimes <- *duration

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
	case err := <- errChan:
		return nil, err
	default:
		return &total, nil
	}
}

func (nw *NodeWorker) executeTask(ctx context.Context, task models.Task) (*time.Duration, error) {
	start := time.Now()
	containerID, err := nw.client.CreateTask(ctx, task.ImageID, task.Config.CPUs)
	if err != nil {
		return nil, fmt.Errorf("create task with id %s: %w", task.ID, err)
	}

	err = nw.client.StartTask(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("start task with id %s: %w", task.ID, err)
	}

	err = nw.client.PauseTask(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("pause task with id %s: %w", task.ID, err)
	}

	end := time.Since(start)

	return &end, nil
}
