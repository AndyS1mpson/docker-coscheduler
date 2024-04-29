package resources

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

type nodeClient interface {
	GetNodeResources(ctx context.Context) (*models.NodeResources, error)
}
