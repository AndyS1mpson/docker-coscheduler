package controllers

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// BuildTask обработчик создания образа для задачи
func (s *Server) BuildTask(ctx context.Context, req *task.BuildTaskRequest) (*task.BuildTaskResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	t, err := s.service.BuildTask(ctx, models.ImageArchive{File: req.ImageArchive}, req.TaskTitle)
	if err != nil {
		return nil, fmt.Errorf("build task: %w", err)
	}

	return &task.BuildTaskResponse{
		TaskId:  t.ID,
		ImageId: t.ImageID,
		Node: &task.Node{
			Host: t.Node.Host,
			Port: t.Node.Port,
		},
		Status: string(t.Status),
	}, nil
}
