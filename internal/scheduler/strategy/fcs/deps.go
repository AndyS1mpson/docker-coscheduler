package fcs

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

// NodeClient представляет систему планирования задач
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
	GetTaskInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error)
}

type taskHub interface {
	ArchiveImageToTar(imageDir string, tarName string) (*models.ImageArchive, error)
}
