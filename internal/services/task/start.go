package task

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// StartTask запуск задачи на ноде
func (s *Service) StartTask(ctx context.Context, containerID string) error {
	info, err := s.dockerClient.GetContainerInfo(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get container info: %w", err)
	}

	if info.State == models.ContainerStateRunning {
		return fmt.Errorf("container is already running")
	}

	return s.dockerClient.StartContainer(ctx, containerID)
}
