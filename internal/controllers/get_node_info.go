package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetNodeInfo получение информации о ноде
func (s *Server) GetNodeInfo(ctx context.Context, req *emptypb.Empty) (*task.GetNodeInfoResponse, error) {
	info := s.service.GetNodeInfo(ctx)

	return &task.GetNodeInfoResponse{
		CpuNums: info.CPUNums,
		Uri:     info.URI,
		Port:    info.Port,
	}, nil
}
