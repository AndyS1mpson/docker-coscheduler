package strategy

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

// RoundRobinStrategy представляет стратегию выполнения задач, в которой назначение задачи узлу идет по алгоритму round-robin.
// Как только задача на каком-либо узле завершает выполнение, на узле запускается следующая задача из списка.
// Если список пуст, на узле ничего больше не запускается
type RoundRobinStrategy[T nodeClient] struct {
	nodes   map[models.Node]T
	taskHub taskHub
	delay   time.Duration
}

// NewRoundRobinStrategy конструктор создания RoundRobinStrategy
func NewRoundRobinStrategy[T nodeClient](nodes map[models.Node]T, taskHub taskHub, taskDelay time.Duration) *RoundRobinStrategy[T] {
	return &RoundRobinStrategy[T]{
		nodes:   nodes,
		taskHub: taskHub,
		delay:   taskDelay,
	}
}

// Execute выполняет стратегию на указанных узлах с задачами
func (s *RoundRobinStrategy[T]) Execute(ctx context.Context, tasks []models.StrategyTask) (time.Duration, error) {
	nodesInfo := slices.Keys(s.nodes)

	// мапа, содержащая информацию о том, свободна ли нода или на ней выполняется задача
	availableNodes := make(map[models.Node]chan struct{}, len(s.nodes))

	for nodeInfo := range s.nodes {
		availableNodes[nodeInfo] = make(chan struct{})
	}

	defer func() {
		for _, ch := range availableNodes {
			close(ch)
		}
	}()

	g, ctx := errgroup.WithContext(ctx)

	start := time.Now()

	for idx, task := range tasks {
		g.Go(func() error {
			info := nodesInfo[idx%len(s.nodes)]
			availableNodes[info] <- struct{}{} // ждем пока нода занята выполнением задачи и занимаем ее
			err := s.executeTask(ctx, s.nodes[info], task)
			<-availableNodes[info] // освобождаем ноду

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func (s *RoundRobinStrategy[T]) executeTask(ctx context.Context, node T, task models.StrategyTask) error {
	archive, err := s.taskHub.ArchiveImageToTar(task.FolderName, task.Name)
	if err != nil {
		return fmt.Errorf("archive task to tar: %w", err)
	}

	taskImageID, err := node.BuildTask(ctx, *archive, task.Name)
	if err != nil {
		return fmt.Errorf("build task: %w", err)
	}

	taskID, err := node.CreateTask(ctx, taskImageID, nil)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	err = node.StartTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("start task: %w", err)
	}

	err = node.WaitForTask(ctx, taskID, s.delay)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("task %s executed", taskID), log.Data{})

	return nil
}
