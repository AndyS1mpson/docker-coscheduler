package fcs

import (
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

type FCSStrategy[T nodeClient] struct {
	nodes   []T
	taskHub taskHub
	delay   time.Duration
}

// NewFCSStrategy конструктор создания FCSStrategy
func NewFCSStrategy[T nodeClient](nodes []T, taskHub taskHub, taskDelay time.Duration) *FCSStrategy[T] {
	return &FCSStrategy[T]{
		nodes:   nodes,
		taskHub: taskHub,
		delay:   taskDelay,
	}
}
