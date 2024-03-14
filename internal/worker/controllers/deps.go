package controllers

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

type service interface {
	GetNodeInfo(ctx context.Context) models.Node
	BuildTask(ctx context.Context, imageArchive models.ImageArchive, taskTitle string) (*models.Task, error)
}
