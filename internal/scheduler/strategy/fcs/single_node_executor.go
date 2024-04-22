package fcs

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/scheduler/services/combination"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
	"golang.org/x/sync/errgroup"
)

// SingleNodeExecutor выполняет FCS стратегию на одной ноде и находит оптимальные комбинации
type SingleNodeFCSExecutor[T nodeClient] struct {
	nodeClient         T
	runner             CombinationRunner
	taskCombinationNum int64
	delay              time.Duration
}

// NewSingleNodeFCSExecutor конструктор для SingleNodeFCSExecutor
func NewSingleNodeFCSExecutor[T nodeClient](
	nodeInfo models.Node,
	nodeClient T,
	taskCombinationNum int64,
	measurementTime time.Duration,
	delayTaskFinishedWaiting time.Duration,
) *SingleNodeFCSExecutor[T] {
	runner := NewCombinationRunner(nodeInfo, nodeClient, taskCombinationNum, measurementTime)

	return &SingleNodeFCSExecutor[T]{
		nodeClient:         nodeClient,
		taskCombinationNum: taskCombinationNum,
		runner:             *runner,
		delay:              delayTaskFinishedWaiting,
	}
}

// Execute выполняет поиск оптимальных комбинаций на ноде и последовательный запуск комбинаций из множества оптимальных
func (sn *SingleNodeFCSExecutor[T]) Execute(ctx context.Context, tasks []models.Task) (*time.Duration, error) {
	combinations := combination.GetCombinations(tasks, sn.taskCombinationNum)

	for idx, combo := range combinations {
		duration, err := sn.runner.Run(ctx, combo)
		if err != nil {
			return nil, err
		}

		combo.TotalTime = *duration

		combinations[idx] = combo
	}

	optimalCombinations := sn.findOptimalCombinations(ctx, slices.Values(combinations))

	start := time.Now()
	for _, combo := range optimalCombinations {
		err := sn.runCombination(ctx, combo)
		if err != nil {
			return nil, err
		}
	}

	duration := time.Since(start)

	return &duration, nil
}

func (sn *SingleNodeFCSExecutor[T]) runCombination(ctx context.Context, combination combination.Combination) error {
	g, groupCtx := errgroup.WithContext(ctx)

	taskIDs := make([]string, 0, len(combination.Tasks))

	for task := range combination.Tasks {
		taskIDs = append(taskIDs, task.Config.ContainerID)

		g.Go(func() error {
			return sn.nodeClient.ResumeTask(groupCtx, task.Config.ContainerID)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	err := sn.waitForTasks(ctx, taskIDs)
	if err != nil {
		return err
	}

	return nil
}

func (sn *SingleNodeFCSExecutor[T]) findOptimalCombinations(ctx context.Context, combinations []combination.Combination) []combination.Combination {
	resultCombinations := make([]combination.Combination, 0, len(combinations))

	sort.Slice(combinations, func(i, j int) bool { return combinations[i].TotalTime < combinations[j].TotalTime })

	tasks := make(map[models.Task]struct{})

	isUnique := true

	for _, combo := range combinations {
		for task := range combo.Tasks {
			if _, ok := tasks[task]; ok {
				isUnique = false

				break
			}
		}
		if isUnique {
			for task := range combo.Tasks {
				tasks[task] = struct{}{}
			}

			resultCombinations = append(resultCombinations, combo)
		}
	}

	return resultCombinations
}

func (sn *SingleNodeFCSExecutor[T]) waitForTasks(ctx context.Context, taskIDs []string) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, taskID := range taskIDs {
		g.Go(func() error {
			isRunning := true

			for isRunning {
				info, err := sn.nodeClient.GetTaskInfo(ctx, taskID)
				if err != nil {
					return fmt.Errorf("get task info: %w", err)
				}

				if info.State == models.ContainerStateRunning {
					time.Sleep(sn.delay)
					continue
				} else if info.State == models.ContainerStateExited && info.ExitCode == 0 {
					isRunning = false
				} else {
					return fmt.Errorf("task crashed with status: %w and exit code: %d", info.State, info.ExitCode)
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
