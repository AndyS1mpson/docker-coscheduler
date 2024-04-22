package fcs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

// FCSStrategy (Fastest Combination Speed) представляет стратегию выполнения задач,
// при которой задачи предварительно раскидываются на доступные узлы и билдятся образы,
// после чего на каждом узле запускаются всевозможные комбинации задач (количество задач в комбинации taskCombinationNum)
// и происходит поиск множества самых "быстрых" комбинаций, которые в последствии и выполняются на узлах
type FCSStrategy[T nodeClient] struct {
	nodes              map[models.Node]T
	runners            map[models.Node]SingleNodeFCSExecutor[T]
	taskHub            taskHub
	taskCombinationNum int64
	delay              time.Duration
}

// NewFCSStrategy конструктор создания FCSStrategy
func NewFCSStrategy[T nodeClient](
	nodes map[models.Node]T,
	taskHub taskHub,
	taskDelay time.Duration,
	taskCombinationNum int64,
	measurementTime time.Duration,
) *FCSStrategy[T] {
	runners := make(map[models.Node]SingleNodeFCSExecutor[T], 0)

	for nodeInfo, node := range nodes {
		runners[nodeInfo] = *NewSingleNodeFCSExecutor(nodeInfo, node, taskCombinationNum, measurementTime, taskDelay)
	}

	return &FCSStrategy[T]{
		nodes:              nodes,
		runners:            runners,
		taskHub:            taskHub,
		taskCombinationNum: taskCombinationNum,
		delay:              taskDelay,
	}
}

// Execute выполняет задачи на узлах по FCS стратегии
func (f *FCSStrategy[T]) Execute(ctx context.Context, tasks []models.StrategyTask) (*time.Duration, error) {
	buildedTasks, err := f.buildTasksOnNodes(ctx, tasks)
	if err != nil {
		log.Error(err, log.Data{})

		return nil, err
	}

	var totalTime time.Duration

	for nodeInfo, tasks := range buildedTasks {
		nodeRunner := f.runners[nodeInfo]
		duration, err := nodeRunner.Execute(ctx, tasks)
		if err != nil {
			return nil, err
		}

		totalTime += *duration
	}

	return &totalTime, nil
}

// buildTasksOnNodes раскидывает задачи по нодам и билдит их
func (f *FCSStrategy[T]) buildTasksOnNodes(ctx context.Context, tasks []models.StrategyTask) (map[models.Node][]models.Task, error) {
	buildedTasks := make(map[models.Node][]models.Task, 0)
	for node := range buildedTasks {
		buildedTasks[node] = make([]models.Task, 0)
	}

	nodes := slices.Keys(f.nodes)

	g, ctx := errgroup.WithContext(ctx)

	var mu sync.Mutex

	for idx, task := range tasks {
		g.Go(func() error {
			nodeInfo := nodes[idx%len(f.nodes)]

			buildedTask, err := f.buildTask(ctx, task, nodeInfo)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()

			buildedTasks[nodeInfo] = append(buildedTasks[nodeInfo], *buildedTask)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return buildedTasks, nil
}

// buildTask сборка задачи на ноде
func (f *FCSStrategy[T]) buildTask(ctx context.Context, task models.StrategyTask, nodeInfo models.Node) (*models.Task, error) {
	buildedTask := models.Task{
		ID:    uuid.NewString(),
		Title: task.Name,
	}

	taskArchive, err := f.taskHub.ArchiveImageToTar(task.FolderName, task.Name)
	if err != nil {
		return nil, fmt.Errorf("archive task: %w", err)
	}

	imageID, err := f.nodes[nodeInfo].BuildTask(ctx, *taskArchive, task.Name)
	if err != nil {
		return nil, fmt.Errorf("build task: %w", err)
	}

	buildedTask.ImageID = imageID
	buildedTask.Node = &nodeInfo
	buildedTask.Status = models.TaskStatusBuild

	return &buildedTask, nil
}
