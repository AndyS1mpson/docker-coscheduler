package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetNodeInfo получение информации о ноде
func (c *Client) GetNodeInfo(ctx context.Context) (*models.Node, error) {
	info, err := c.externalClient.GetNodeInfo(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return &models.Node{
		Port:    info.Port,
		Host:    info.Uri,
		CPUNums: info.CpuNums,
	}, nil
}
