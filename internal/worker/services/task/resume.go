package task

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// ResumeTask продолжение выполнения задачи
func (s *Service) ResumeTask(ctx context.Context, containerID string) error {
	info, err := s.dockerClient.GetContainerInfo(ctx, containerID)
	if err != nil {
		return fmt.Errorf("get container info: %w", err)
	}

	if info.State != models.ContainerStatePaused {
		if info.State == models.ContainerStateRunning {
			return fmt.Errorf("container is already unpaused")
		}

		return fmt.Errorf("container is not paused")
	}

	return s.dockerClient.UnpauseContainer(ctx, containerID)
}
