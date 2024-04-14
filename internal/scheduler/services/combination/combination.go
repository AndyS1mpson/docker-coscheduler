package combination

import (
	"time"

	"github.com/google/uuid"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

// Combination комбинация задач для FCS стратегии
type Combination struct {
	Title     string
	Tasks     map[models.Task]struct{}
	TotalTime time.Duration
}

// GetCombinations получает всевозможные комбинации задач размером taskCombintationNum
func GetCombinations(tasks []models.Task, taskCombintationNum int64) map[string]Combination {
	combinations := make(map[string]Combination)
	generateHelper(tasks, taskCombintationNum, 0, []models.Task{}, combinations)

	return combinations
}

func generateHelper(tasks []models.Task, taskCombintationNum int64, index int, currentCombination []models.Task, combinations map[string]Combination) {
	combinationID := uuid.NewString()

	if len(currentCombination) == int(taskCombintationNum) {
		combinations[combinationID] = Combination{
			Title: combinationID,
			Tasks: slices.ToMap(currentCombination),
		}

		return
	}

	if index >= len(tasks) {
		return
	}

	task := tasks[index]
	generateHelper(tasks, taskCombintationNum, index+1, append(currentCombination, task), combinations)
	generateHelper(tasks, taskCombintationNum, index+1, currentCombination, combinations)
}
