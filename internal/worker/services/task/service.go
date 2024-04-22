package task

import (
	"context"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Service содержит логику работы с тасками
type Service struct {
	dockerClient dockerClient
	unpacker     unpacker
	measurer     measurer
	config       models.Node
}

// NewService конструктор для Service
func NewService(dockerClient dockerClient, unpacker unpacker, measurer measurer, nodeURI string, nodePort int64, cpuNums int64) *Service {
	return &Service{
		dockerClient: dockerClient,
		unpacker:     unpacker,
		measurer:     measurer,
		config: models.Node{
			Host:    nodeURI,
			Port:    nodePort,
			CPUNums: cpuNums,
		},
	}
}

// GetNodeInfo получение информации о ноде на которой работает воркер
func (s *Service) GetNodeInfo(ctx context.Context) models.Node {
	return s.config
}

// GetContainerInfo получение информации о контейнере с задачей
func (s *Service) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	return s.dockerClient.GetContainerInfo(ctx, containerID)
}

// MeasureTaskSpeed замеряет время выполнения задачи на ноде
func (s *Service) MeasureTaskSpeed(ctx context.Context, task models.Task, duration time.Duration) (time.Duration, error) {
	return s.measurer.Measure(ctx, task, duration)
}
