package controllers

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// ResumeTask обработчик возобновления задачи
func (s *Server) ResumeTask(ctx context.Context, req *task.ResumeTaskRequest) (*emptypb.Empty, error) {
	err := s.service.ResumeTask(ctx, req.ContainerId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
