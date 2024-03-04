package task

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// PauseTask остановка выполнения задачи
func (s *Service) PauseTask(ctx context.Context, containerID string) error {
	info, err := s.dockerClient.GetContainerInfo(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get container info: %w", err)
	}

	if info.State != models.ContainerStateRunning {
		if info.State == models.ContainerStatePaused {
			return fmt.Errorf("container is already paused")
		}

		return fmt.Errorf("container is not running")
	}

	return s.dockerClient.PauseContainer(ctx, containerID)
}
