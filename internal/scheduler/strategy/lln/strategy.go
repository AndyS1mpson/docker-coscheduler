package lln

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
)

// LeastLoadedNode представляет стратегию выполнения задач, в которой задачи назначаются самому незагруженному узлу.
type LeastLoadedNode[T nodeClient] struct {
	storage            storage
	repository         repository
	nodeResourcesCache nodeResourcesCache
	strategiesCache    strategiesCache

	nodes                   map[models.Node]T
	taskHub                 taskHub
	delay                   time.Duration
	maxNodeCPUUsageLimit    *float64
	maxNodeMemoryUsageLimit *float64
}

// NewLeastLoadedNode конструктор создания LeastLoadedNode
func NewLeastLoadedNode[T nodeClient](
	storage storage,
	repository repository,
	nodeResourcesCache nodeResourcesCache,
	strategiesCache strategiesCache,
	nodes map[models.Node]T,
	taskHub taskHub,
	taskDelay time.Duration,
	maxNodeCPUUsageLimit *float64,
	maxNodeMemoryUsageLimit *float64,
) *LeastLoadedNode[T] {
	return &LeastLoadedNode[T]{
		storage:                 storage,
		repository:              repository,
		nodeResourcesCache:      nodeResourcesCache,
		strategiesCache:         strategiesCache,
		nodes:                   nodes,
		taskHub:                 taskHub,
		delay:                   taskDelay,
		maxNodeCPUUsageLimit:    maxNodeCPUUsageLimit,
		maxNodeMemoryUsageLimit: maxNodeMemoryUsageLimit,
	}
}

// Execute выполняет стратегию на указанных узлах с задачами
func (l *LeastLoadedNode[T]) Execute(ctx context.Context, experimentID uuid.UUID, tasks []models.StrategyTask) (time.Duration, error) {

	g, ctx := errgroup.WithContext(ctx)

	start := time.Now()

	for _, task := range tasks {
		g.Go(func() error {
			for {
				sortedNodes := l.nodeResourcesCache.SortNodesByResources(ctx, l.maxNodeCPUUsageLimit, l.maxNodeMemoryUsageLimit)

				if len(sortedNodes) != 0 {
					_, err := l.executeTask(ctx, l.nodes[sortedNodes[0]], task)
					if err != nil {
						return err
					}
				}
			}
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	duration := time.Since(start)

	err := l.saveExperimentResults(ctx, experimentID, tasks, duration)
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (l *LeastLoadedNode[T]) executeTask(ctx context.Context, node T, task models.StrategyTask) (time.Duration, error) {
	archive, err := l.taskHub.ArchiveImageToTar(task.FolderName, task.Name)
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

	err = node.WaitForTask(ctx, taskID, l.delay)
	if err != nil {
		return 0, err
	}

	log.Info(fmt.Sprintf("task %s executed", taskID), log.Data{})

	return time.Since(start), nil
}

func (l *LeastLoadedNode[T]) saveExperimentResults(
	ctx context.Context,
	experimentID uuid.UUID,
	tasks []models.StrategyTask,
	totalTime time.Duration,
) error {
	return l.storage.Tx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		strategyID, err := l.getStrategyID(ctx, tx, models.StrategyNameRoundRobin)
		if err != nil {
			return fmt.Errorf("get strategy id: %w", err)
		}

		id, err := l.repository.SaveExperimentResultTx(ctx, tx, models.ExperimentResult{
			IdempotencyKey: experimentID.String(),
			StrategyID:     strategyID,
			ExecutionTime:  totalTime,
		})
		if err != nil {
			return fmt.Errorf("save experiment result: %w", err)
		}

		_, err = l.repository.SaveExperimentStrategyTasksTx(ctx, tx, id, tasks)
		if err != nil {
			return fmt.Errorf("save experiment tasks info: %w", err)
		}

		return nil
	})
}

// getStrategyID получает или создает стратегию для задачи в базе
func (l *LeastLoadedNode[T]) getStrategyID(ctx context.Context, tx *sqlx.Tx, strategyName models.StrategyName) (int64, error) {
	strategy := l.strategiesCache.GetByName(ctx, models.StrategyNameRoundRobin)

	if strategy == nil {
		return l.repository.CreateStrategy(ctx, models.Strategy{Name: strategyName})
	}

	return strategy.ID, nil
}
