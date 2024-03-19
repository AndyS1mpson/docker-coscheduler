package task

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/google/uuid"
)

// CreateTask создание docker контейнера для задачи
func (s *Service) CreateTask(ctx context.Context, task models.Task, cpuOpt *models.CPUSet) (*models.Task, error) {
	containerID, err := s.dockerClient.CreateContainer(ctx, task.ImageID, cpuOpt, uuid.NewString())
	if err != nil {
		return nil, fmt.Errorf("create docker container for task: %w", err)
	}

	task.Status = models.TaskStatusCreated
	task.Config.ContainerID = containerID

	return &task, nil
}
