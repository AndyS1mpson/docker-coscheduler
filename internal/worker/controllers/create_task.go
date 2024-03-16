package controllers

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// CreateTask обработчик создания контейнера с задачей 
func (s *Server) CreateTask(ctx context.Context, req *task.CreateTaskRequest) (*task.CreateTaskResponse, error) {
	createdTask, err := s.service.CreateTask(ctx, models.Task{ImageID: req.ImageId}, models.CPUSet{From: req.CpusOpt.From, Count: req.CpusOpt.Count})
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return &task.CreateTaskResponse{
		ContainerId: createdTask.Config.ContainerID,
		Status:      string(createdTask.Status),
	}, nil
}
