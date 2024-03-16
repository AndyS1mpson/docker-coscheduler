package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PauseTask обработчик остановки выполнения задачи
func (s *Server) PauseTask(ctx context.Context, req *task.PauseTaskRequest) (*emptypb.Empty, error) {
	err := s.service.PauseTask(ctx, req.ContainerId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
