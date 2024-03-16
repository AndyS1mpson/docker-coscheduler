package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"

	"google.golang.org/protobuf/types/known/emptypb"
)

// StartTask обработчик запуска задачи
func (s *Server) StartTask(ctx context.Context, req *task.StartTaskRequest) (*emptypb.Empty, error) {
	err := s.service.StartTask(ctx, req.ContainerId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
