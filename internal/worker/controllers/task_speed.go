package controllers

import (
	"context"
	"fmt"

	"github.com/AndyS1mpson/docker-coscheduler/generated/task"
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"google.golang.org/protobuf/types/known/durationpb"
)

// MeasureTaskSpeed обработчик измерения времени выполнения задачи
func (s *Server) MeasureTaskSpeed(ctx context.Context, req *task.MeasureTaskSpeedRequest) (*task.MeasureTaskSpeedResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	speed, err := s.service.MeasureTaskSpeed(ctx, models.Task{
		Config: &models.Config{
			ContainerID: req.ContainerId,
			CPUs: &models.CPUSet{
				From: req.CpusOpt.From,
				Count: req.CpusOpt.Count,
			},
		}},
		req.Duration.AsDuration(),
	)
	if err != nil {
		return nil, fmt.Errorf("build task: %w", err)
	}

	return &task.MeasureTaskSpeedResponse{
		Time: durationpb.New(speed),
	}, nil
}
