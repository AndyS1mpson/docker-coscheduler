package controllers

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/worker/services/resourcer"
)

// GetNodeResources получение информации о занятых ресурсах ноды
func (s *Server) GetNodeResources(ctx context.Context, req *emptypb.Empty) (*task.GetNodeResourcesResponse, error) {
	resources, err := resourcer.GetNodeResources()
	if err != nil {
		return nil, err
	}

	return &task.GetNodeResourcesResponse{
		Cpu:    resources.CPUUtilization,
		Memory: resources.MemoryUtilization,
	}, nil
}
