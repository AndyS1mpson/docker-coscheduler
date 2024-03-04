package strategy

import (
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Config конфигурация для стратегий
type Config struct {
	NodesUri         []string
	Tasks            []models.Task
	WaitForTaskDelay time.Duration
}
