package task

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

type dockerClient interface {
	GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error)
	BuildImage(ctx context.Context, dirName string) (string, error)
	CreateContainer(ctx context.Context, imageID string, cpuSet models.CPUSet, containerName string) (string, error)
	StartContainer(ctx context.Context, containerID string) error
	PauseContainer(ctx context.Context, containerID string) error
	UnpauseContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string) error
}

type unpacker interface {
	UnpackTarArchive(archiveFile models.ImageArchive, dirName string) (string, error)
}
