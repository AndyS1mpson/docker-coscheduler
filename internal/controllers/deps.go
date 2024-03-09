package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/services/task"
)

type service interface {
	GetNodeInfo(ctx context.Context) task.Config 
}
