package cache

import (
	"sort"
	"sync"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// Cache кэш, хранящий времена выполнения последних N выполнившихся на ноде задач
type Cache struct {
	mu   sync.RWMutex
	data map[models.Node][]time.Duration
	num  int64
}

// NewCache конструктор для Cache
func NewCache(nodes []models.Node, num int64) *Cache {
	data := make(map[models.Node][]time.Duration, 0)
	for _, node := range nodes {
		data[node] = make([]time.Duration, num)
	}

	return &Cache{
		data: data,
		num:  num,
	}
}

// SetExecutionTime сохраняет новое время выолнения задачи на ноде в кэш
func (c *Cache) SetExecutionTime(node models.Node, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.data[node]; !ok {
		c.data[node] = make([]time.Duration, c.num)
	}

	c.data[node] = append(c.data[node], duration)
	if len(c.data[node]) > int(c.num) {
		c.data[node] = c.data[node][1:]
	}
}

// SortedNodesByAvg сортировка нод по возрастанию среднего времени выполнения задач
func (c *Cache) SortedNodesByAvg() []models.Node {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	type NodeAvg struct {
		Node models.Node
		Avg  time.Duration
	}

	nodeAvgs := make([]NodeAvg, 0, len(c.data))

	for node, durations := range c.data {
		sum := time.Duration(0)
		for _, duration := range durations {
			sum += duration
		}

		average := sum / time.Duration(len(durations))

		nodeAvgs = append(nodeAvgs, NodeAvg{Node: node, Avg: average})
	}

	sort.Slice(nodeAvgs, func(i, j int) bool {
		return nodeAvgs[i].Avg < nodeAvgs[j].Avg
	})

	sortedNodes := make([]models.Node, len(nodeAvgs))
	for i, na := range nodeAvgs {
		sortedNodes[i] = na.Node
	}

	return sortedNodes
}
