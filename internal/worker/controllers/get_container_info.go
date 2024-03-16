package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
)

// GetContainerInfo обработчик получения информации о состоянии контейнера с задачей
func (s *Server) GetContainerInfo(ctx context.Context, req *task.GetContainerInfoRequest) (*task.GetContainerInfoResponse, error) {
	info, err := s.service.GetContainerInfo(ctx, req.ContainerId)
	if err != nil {
		return nil, err
	}

	return &task.GetContainerInfoResponse{
		Id:       info.ID,
		State:    string(info.State),
		ExitCode: info.ExitCode,
	}, nil
}
