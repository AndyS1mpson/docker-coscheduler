package worker

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"google.golang.org/protobuf/types/known/emptypb"
)

// GetNodeResources получение информации о загруженности ноды
func (c *Client) GetNodeResources(ctx context.Context) (*models.NodeResources, error) {
	info, err := c.externalClient.GetNodeResources(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return &models.NodeResources{
		CPUUtilization:    info.Cpu,
		MemoryUtilization: info.Memory,
	}, nil
}
