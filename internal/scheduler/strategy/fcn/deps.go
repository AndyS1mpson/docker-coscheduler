package fcn

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// nodeClient представляет систему планирования задач
type nodeClient interface {
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

type cache interface {
	SetExecutionTime(node models.Node, duration time.Duration)
	SortedNodesByAvg() []models.Node
}

type taskHub interface {
	ArchiveImageToTar(imageDir string, tarName string) (*models.ImageArchive, error)
}