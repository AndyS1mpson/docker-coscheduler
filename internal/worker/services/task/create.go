package task

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// CreateTask создание docker контейнера для задачи
func (s *Service) CreateTask(ctx context.Context, task models.Task, cpuOpt models.CPUSet) (*models.Task, error) {
	containerID, err := s.dockerClient.CreateContainer(ctx, task.ImageID, cpuOpt, task.ImageID)
	if err != nil {
		return nil, fmt.Errorf("create docker container for task: %w", err)
	}

	task.Status = models.TaskStatusCreated
	task.Config.ContainerID = containerID

	return &task, nil
}
