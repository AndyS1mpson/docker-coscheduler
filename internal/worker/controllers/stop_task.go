package controllers

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// StopTask обработчик остановки задачи
func (s *Server) StopTask(ctx context.Context, req *task.StopTaskRequest) (*emptypb.Empty, error) {
	err := s.service.StopTask(ctx, req.ContainerId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
