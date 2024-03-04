package strategy

// import (
// 	"context"

// 	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
// )

// // SequentialStrategy представляет последовательную стратегию выполнения задач.
// // Последовательная стратегия - запуск всех задач по очереди монопольно на каждом узле.
// // Как только задача на каком-либо узле завершает выполнение, на узле запускается следующая задача из списка.
// // Если список пуст, на узле ничего больше не запускается
// type SequentialStrategy struct {
// 	client NodeClient
// }

// // NewSequentialStrategy конструктор создания SequentialStrategy
// func NewSequentialStrategy(client NodeClient) *SequentialStrategy {
// 	return &SequentialStrategy{
// 		client: client,
// 	}
// }

// // Execute выполняет стратегию на указанных узлах с задачами
// func (s *SequentialStrategy) Execute(ctx context.Context, nodes []models.Node, tasks []models.Task) {
// 	tasksRef := make(chan models.Task, len(tasks))
	
// 	for _, task := range tasks {
// 		tasksRef <- task
// 	}
	
// 	for _, node := range nodes {
// 		go func(node models.Node) {
// 			for {
// 				select {
// 				case task := <-tasksRef:
// 					s.client.BuildTask(ctx, node, task)
// 				}
// 			}
// 		}(node)
// 	}
// }
