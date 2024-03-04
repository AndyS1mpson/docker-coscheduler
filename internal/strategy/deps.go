package strategy

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

//go:generate mockgen -source=deps.go -destination=./mocks/mock.go

// NodeClient представляет систему планирования задач
type NodeClient interface {
	NodeInfo(ctx context.Context, uri string) (*models.Node, error)
	BuildTask(ctx context.Context, node *models.Node, image *models.ImageArchive, taskName string) (*models.Task, error)
	CreateTask(ctx context.Context, task *models.Task, cpuset *models.CPUSet) (*models.Task, error)
	StartTask(ctx context.Context, task *models.Task) (*models.Task, error)
	PauseTask(ctx context.Context, task *models.Task) (*models.Task, error)
	ResumeTask(ctx context.Context, task *models.Task) (*models.Task, error)
	RunTaskFromTuple(ctx context.Context, node *models.Node, taskName string, cpuSet *models.CPUSet) (*models.Task, error)
	StopTask(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitForTask(ctx context.Context, task *models.Task) (int64, error)
	IsRunning(ctx context.Context, task *models.Task) (bool, error)
	InitSession(ctx context.Context, node *models.Node) error
	UpdateCpus(ctx context.Context, task *models.Task, cpuSet *models.CPUSet) (*models.Task, error)
}
