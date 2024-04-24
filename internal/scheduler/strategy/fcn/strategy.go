package fcn

import (
	"context"
	"fmt"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"golang.org/x/sync/errgroup"
)

// FCNStrategy (Fastest computing node) - стратегия,
// при которой приоритет в выполнении задачи отдается самому "мощному" вычислительному узлу из свободных.
// При этом количество одновременно выполняющихся задач вычисляется пропорционально мощности ноды.
type FCNStrategy[T nodeClient] struct {
	nodes               map[models.Node]T
	cache               cache
	taskHub             taskHub
	computingTaskNum    int64
	delay               time.Duration
	taskMeasurementTime time.Duration
}

// NewFCNStrategy конструктор для FCNStrategy
func NewFCNStrategy[T nodeClient](nodes map[models.Node]T, cache cache, taskHub taskHub, computingTaskNum int64, delay time.Duration, taskMeasurementTime time.Duration) *FCNStrategy[T] {
	return &FCNStrategy[T]{
		nodes:               nodes,
		cache:               cache,
		taskHub:             taskHub,
		computingTaskNum:    computingTaskNum,
		delay:               delay,
		taskMeasurementTime: taskMeasurementTime,
	}
}

// Execute выполняет задачи на узлах по FCN стратегии
func (f *FCNStrategy[T]) Execute(ctx context.Context, tasks []models.StrategyTask) (time.Duration, error) {

	availableNodes := make(map[models.Node]chan struct{}) // мапа, которая отвечает за то, что на ноде одновременно выполняется не более computingTaskNum задач
	for node := range f.nodes {
		availableNodes[node] = make(chan struct{}, f.computingTaskNum)
	}

	defer func() {
		for _, ch := range availableNodes {
			close(ch)
		}
	}()

	g, ctx := errgroup.WithContext(ctx)

	start := time.Now()

	for _, task := range tasks {
		g.Go(func() error {
			sortedNodes := f.cache.SortedNodesByAvg()

			for {
				for _, node := range sortedNodes {
					select {
					case availableNodes[node] <- struct{}{}:
						defer func() { <-availableNodes[node] }()

						duration, err := f.executeTask(ctx, f.nodes[node], task)
						if err != nil {
							return err
						}

						f.cache.SetExecutionTime(node, duration)

						return nil
					default:
						continue
					}
				}
			}
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	duration := time.Since(start)

	return duration, nil
}

func (f *FCNStrategy[T]) executeTask(ctx context.Context, node T, task models.StrategyTask) (time.Duration, error) {
	archive, err := f.taskHub.ArchiveImageToTar(task.FolderName, task.Name)
	if err != nil {
		return 0, fmt.Errorf("archive task to tar: %w", err)
	}

	taskImageID, err := node.BuildTask(ctx, *archive, task.Name)
	if err != nil {
		return 0, fmt.Errorf("build task: %w", err)
	}

	start := time.Now()

	taskID, err := node.CreateTask(ctx, taskImageID, nil)
	if err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}

	err = node.StartTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("start task: %w", err)
	}

	err = node.WaitForTask(ctx, taskID, f.delay)
	if err != nil {
		return 0, err
	}

	log.Info(fmt.Sprintf("task %s executed", taskID), log.Data{})

	return time.Since(start), nil
}
