package task

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// BuildTask создание docker образа для задачи на воркере
func (s *Service) BuildTask(ctx context.Context, imageArchive models.ImageArchive, taskTitle string) (*models.Task, error) {
	id := uuid.NewString()

	taskDir, err := s.unpacker.UnpackTarArchive(imageArchive, id)
	if err != nil {
		return nil, fmt.Errorf("unpack archive: %w", err)
	}

	imageID, err := s.dockerClient.BuildImage(ctx, taskDir)
	if err != nil {
		return nil, fmt.Errorf("build docker image: %w", err)
	}

	return &models.Task{
		ID:      id,
		ImageID: imageID,
		Node: &models.Node{
			Host: s.config.URI,
			Port: s.config.Port,
		},
		Status: models.TaskStatusBuild,
		Title: taskTitle,
	}, nil
}
