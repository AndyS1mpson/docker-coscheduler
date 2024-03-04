package task

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Service содержит логику работы с тасками
type Service struct {
	dockerClient dockerClient
	unpacker     unpacker
	config       NodeConfig
}

// NewService конструктор для Service
func NewService(dockerClient dockerClient, uunpacker unpacker, config NodeConfig) *Service {
	return &Service{
		dockerClient: dockerClient,
		unpacker:     uunpacker,
		config:       config,
	}
}

// GetNodeInfo получение информации о ноде на которой работает воркер
func (s *Service) GetNodeInfo(ctx context.Context) NodeConfig {
	return s.config
}

// GetContainerInfo получение информации о контейнере с задачей
func (s *Service) GetContainerInfo(ctx context.Context, containerID string) (*models.ContainerInfo, error) {
	return s.dockerClient.GetContainerInfo(ctx, containerID)
}
