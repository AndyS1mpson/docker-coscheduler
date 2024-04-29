package lln

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

// nodeClient представляет систему планирования задач
type nodeClient interface {
	GetNodeResources(ctx context.Context) (*models.NodeResources, error)
	BuildTask(ctx context.Context, archive models.ImageArchive, taskTitle string) (string, error)
	CreateTask(ctx context.Context, imageID string, cpuSet *models.CPUSet) (string, error)
	StartTask(ctx context.Context, containerID string) error
	PauseTask(ctx context.Context, containerID string) error
	MeasureTaskSpeed(
		ctx context.Context,
		containerID string,
		cpuSet models.CPUSet,
		duration time.Duration,
	) (time.Duration, error)
	ResumeTask(ctx context.Context, containerID string) error
	WaitForTask(ctx context.Context, taskID string, delay time.Duration) error
}

type taskHub interface {
	ArchiveImageToTar(imageDir string, tarName string) (*models.ImageArchive, error)
}

type storage interface {
	Tx(ctx context.Context, f func(context.Context, *sqlx.Tx) error) error
}

type repository interface {
	SaveExperimentResultTx(ctx context.Context, q *sqlx.Tx, result models.ExperimentResult) (int64, error)
	SaveExperimentStrategyTasksTx(
		ctx context.Context,
		q *sqlx.Tx,
		experimentID int64,
		tasks []models.StrategyTask,
	) (int64, error)
	CreateStrategy(ctx context.Context, strategy models.Strategy) (int64, error)
}

type nodeResourcesCache interface {
	SortNodesByResources(ctx context.Context, maxUsageCPULimit *float64, maxUsageMemoryLimit *float64) []models.Node
}

type strategiesCache interface {
	GetByName(ctx context.Context, name models.StrategyName) *models.Strategy
}
