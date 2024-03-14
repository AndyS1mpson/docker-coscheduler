package task

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// StopTask остановка задачи
func (s *Service) StopTask(ctx context.Context, containerID string) error {
	info, err := s.dockerClient.GetContainerInfo(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get container info: %w", err)
	}

	if info.State != models.ContainerStateRunning {
		return fmt.Errorf("container is not running")
	}

	return s.dockerClient.StopContainer(ctx, containerID)
}
