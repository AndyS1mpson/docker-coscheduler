package resources

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
	"golang.org/x/sync/errgroup"
)

// Cache кэш, хранящий состояние загруженности нод
type Cache[T nodeClient] struct {
	mu   sync.RWMutex
	data map[models.Node]models.NodeResources

	nodes          map[models.Node]T
	reloadInterval time.Duration
}

// NewCache конструктор для Cache
func NewCache[T nodeClient](nodes map[models.Node]T, reloadInterval time.Duration) *Cache[T] {
	return &Cache[T]{
		nodes:          nodes,
		reloadInterval: reloadInterval,
		data:           make(map[models.Node]models.NodeResources),
	}
}

// StartLoading запускает асинхронное получение информации о ресурсах нод и загрузку ее в кэш
func (c *Cache[T]) StartLoading(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(c.reloadInterval):
				if err := c.Load(ctx); err != nil {
					log.Error(err, log.Data{})
				}
			}
		}
	}()
}

// Load загрузка стратегий в кэш
func (c *Cache[T]) Load(ctx context.Context) error {
	resources := make(map[models.Node]models.NodeResources, len(c.nodes))

	g, ctx := errgroup.WithContext(ctx)

	for node, client := range c.nodes {
		g.Go(func() error {
			info, err := client.GetNodeResources(ctx)
			if err != nil {
				return err
			}
	
			resources[node] = *info

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = resources

	return nil
}

// GetByName получить информацию о загруженности ноды
func (c *Cache[T]) GetNodeResources(ctx context.Context, node models.Node) *models.NodeResources {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c, ok := c.data[node]; ok {
		return &c
	}

	return nil
}

// SortNodesByResources возвращает массив нод, отсортированный по свободным ресурсам
func (c *Cache[T]) SortNodesByResources(ctx context.Context, maxUsageCPULimit *float64, maxUsageMemoryLimit *float64) []models.Node {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var nodes []models.Node

	for node, resources := range c.data {
		if maxUsageCPULimit != nil {
			if resources.CPUUtilization > *maxUsageCPULimit {
				continue
			}
		}

		if maxUsageMemoryLimit != nil {
			if resources.MemoryUtilization > *maxUsageMemoryLimit {
				continue
			}
		}

		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i int, j int) bool {
		return c.data[nodes[i]].CPUUtilization < c.data[nodes[j]].CPUUtilization
	})

	return nodes
}
