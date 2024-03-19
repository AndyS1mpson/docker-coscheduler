package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateTaskResources обработчик обновления ресурсов, выделенных задаче
func (s *Server) UpdateTaskResources(ctx context.Context, req *task.UpdateTaskResourcesRequest) (*emptypb.Empty, error) {
	err := s.service.UpdateTaskResources(ctx, req.ContainerId, models.CPUSet{
		From:  req.CpusOpt.From,
		Count: req.CpusOpt.Count,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
