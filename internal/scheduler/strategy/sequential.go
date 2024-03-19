package strategy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
)

// SequentialStrategy представляет последовательную стратегию выполнения задач.
// Последовательная стратегия - запуск всех задач по очереди монопольно на каждом узле.
// Как только задача на каком-либо узле завершает выполнение, на узле запускается следующая задача из списка.
// Если список пуст, на узле ничего больше не запускается
type SequentialStrategy[T nodeClient] struct {
	nodes   []T
	taskHub taskHub
	delay   time.Duration
}

// NewSequentialStrategy конструктор создания SequentialStrategy
func NewSequentialStrategy[T nodeClient](nodes []T, taskHub taskHub, taskDelay time.Duration) *SequentialStrategy[T] {
	return &SequentialStrategy[T]{
		nodes:   nodes,
		taskHub: taskHub,
		delay:   taskDelay,
	}
}

// Execute выполняет стратегию на указанных узлах с задачами
func (s *SequentialStrategy[T]) Execute(ctx context.Context, tasks []models.StrategyTask) {
	tasksRef := make(chan models.StrategyTask, len(tasks))

	for _, task := range tasks {
		tasksRef <- task
	}

	var wg sync.WaitGroup

	wg.Add(len(s.nodes))

	for _, node := range s.nodes {
		go func(node T) {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasksRef:
					if !ok {
						break
					}

					err := s.executeTask(ctx, node, task)
					if err != nil {
						log.Error(err, log.Data{})
					}
				}
			}
		}(node)
	}

	close(tasksRef)

	wg.Wait()
}

func (s *SequentialStrategy[T]) executeTask(ctx context.Context, node T, task models.StrategyTask) error {
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

	err = s.waitForTask(ctx, node, taskID)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("task %s executed", taskID), log.Data{})

	return nil
}

func (s *SequentialStrategy[T]) waitForTask(ctx context.Context, node T, taskID string) error {
	isRunning := true

	for isRunning {
		info, err := node.GetTaskInfo(ctx, taskID)
		if err != nil {
			return fmt.Errorf("get task info: %w", err)
		}

		if info.State == models.ContainerStateRunning {
			time.Sleep(s.delay)
			continue
		} else if info.State == models.ContainerStateExited && info.ExitCode == 0 {
			isRunning = false
		} else {
			return fmt.Errorf("task crashed with status: %w and exit code: %d", info.State, info.ExitCode)
		}
	}

	return nil
}
