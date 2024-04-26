package experiment

import "github.com/AndyS1mpson/docker-coscheduler/internal/models"

// SearchCriteria критерии для поиска результатов экспериментов
type SearchCriteria struct {
	IDs             []int64
	IdempotencyKeys []string
	StrategyNames   []models.StrategyName
}
