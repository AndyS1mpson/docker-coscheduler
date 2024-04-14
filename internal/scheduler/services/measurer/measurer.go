package measurer

import (
	"context"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

type SpeedMeasurer struct {

}

// NewSpeedMeasurer конструктор для SpeedMeasurer
func NewSpeedMeasurer() *SpeedMeasurer {
	return &SpeedMeasurer{}
}

// Measure измерение времени выполнения комбинации задач
func (sm *SpeedMeasurer) Measure(ctx context.Context, combination models.Combination) {
	
}
