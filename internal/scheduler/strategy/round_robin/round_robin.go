package round_robin

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

// RoundRobinStrategy представляет стратегию выполнения задач, в которой назначение задачи узлу идет по алгоритму round-robin.
// Как только задача на каком-либо узле завершает выполнение, на узле запускается следующая задача из списка.
// Если список пуст, на узле ничего больше не запускается
type RoundRobinStrategy[T nodeClient] struct {
	storage    storage
	repository repository
	nodes      map[models.Node]T
	taskHub    taskHub
	delay      time.Duration
}

// NewRoundRobinStrategy конструктор создания RoundRobinStrategy
func NewRoundRobinStrategy[T nodeClient](
	storage storage,
	repository repository,
	nodes map[models.Node]T,
	taskHub taskHub,
	taskDelay time.Duration,
) *RoundRobinStrategy[T] {
	return &RoundRobinStrategy[T]{
		storage:    storage,
		repository: repository,
		nodes:      nodes,
		taskHub:    taskHub,
		delay:      taskDelay,
	}
}

// Execute выполняет стратегию на указанных узлах с задачами
func (s *RoundRobinStrategy[T]) Execute(ctx context.Context, experimentID uuid.UUID, tasks []models.StrategyTask) (time.Duration, error) {
	nodesInfo := slices.Keys(s.nodes)

	// мапа, содержащая информацию о том, свободна ли нода или на ней выполняется задача
	availableNodes := make(map[models.Node]chan struct{}, len(s.nodes))

	for nodeInfo := range s.nodes {
		availableNodes[nodeInfo] = make(chan struct{}, 1)
	}

	defer func() {
		for _, ch := range availableNodes {
			close(ch)
		}
	}()

	g, groupCtx := errgroup.WithContext(ctx)

	start := time.Now()

	for idx, task := range tasks {
		g.Go(func() error {
			info := nodesInfo[idx%len(s.nodes)]
			availableNodes[info] <- struct{}{} // ждем пока нода занята выполнением задачи и занимаем ее
			err := s.executeTask(groupCtx, s.nodes[info], task)
			<-availableNodes[info] // освобождаем ноду

			return err
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	end := time.Since(start)

	err := s.saveExperimentResults(ctx, experimentID, tasks, end)
	if err != nil {
		return 0, err
	}

	return end, nil
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

func (s *RoundRobinStrategy[T]) saveExperimentResults(
	ctx context.Context,
	experimentID uuid.UUID,
	tasks []models.StrategyTask,
	totalTime time.Duration,
) error {
	return s.storage.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		id, err := s.repository.SaveExperimentResultTx(ctx, tx, models.ExperimentResult{
			IdempotencyKey: experimentID.String(),
			StrategyName:   models.StrategyNameRoundRobin,
			ExecutionTime:  totalTime,
		})
		if err != nil {
			return fmt.Errorf("save experiment result: %w", err)
		}

		_, err = s.repository.SaveExperimentStrategyTasksTx(ctx, tx, id, tasks)
		if err != nil {
			return fmt.Errorf("save experiment tasks info: %w", err)
		}

		return nil
	})
}
