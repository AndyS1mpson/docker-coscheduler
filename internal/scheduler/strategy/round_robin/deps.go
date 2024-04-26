package round_robin

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

// NodeClient представляет систему планирования задач
type nodeClient interface {
	BuildTask(ctx context.Context, archive models.ImageArchive, taskTitle string) (string, error)
	CreateTask(ctx context.Context, imageID string, cpuSet *models.CPUSet) (string, error)
	StartTask(ctx context.Context, containerID string) error
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
}
