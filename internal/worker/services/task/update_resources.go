package task

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// UpdateTaskResources обновление ресурсов, выделенных под выполнение задачи
func (s *Service) UpdateTaskResources(ctx context.Context, containerID string, cpuSet models.CPUSet) error {
	return s.dockerClient.UpdateContainer(ctx, containerID, cpuSet)
}
