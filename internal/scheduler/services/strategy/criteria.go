package strategy

import (
	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/sql"
)

// SearchCriteria фильтры для поиска стратегий
type SearchCriteria struct {
	IDs   []int64
	Names []models.StrategyName
	PageQuery *sql.PageQuery
}
