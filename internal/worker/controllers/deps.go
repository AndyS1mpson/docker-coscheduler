package controllers

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

type service interface {
	GetNodeInfo(ctx context.Context) models.Node
	BuildTask(ctx context.Context, imageArchive models.ImageArchive, taskTitle string) (*models.Task, error)
	CreateTask(ctx context.Context, task models.Task, cpuOpt *models.CPUSet) (*models.Task, error)
	StartTask(ctx context.Context, containerID string) error
	PauseTask(ctx context.Context, containerID string) error
	ResumeTask(ctx context.Context, containerID string) error
	StopTask(ctx context.Context, containerID string) error
	UpdateTaskResources(ctx context.Context, containerID string, cpuSet models.CPUSet) error
	GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error)
	MeasureTaskSpeed(ctx context.Context, task models.Task, duration time.Duration) (time.Duration, error)
}
