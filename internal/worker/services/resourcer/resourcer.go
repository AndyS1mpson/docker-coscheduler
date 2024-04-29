package resourcer

import (
	"fmt"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// GetNodeResources получение информации о загруженности процессора и памяти узла
func GetNodeResources() (*models.NodeResources, error) {
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("get cpu utilization info: %w", err)
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("get memory utilization info: %w", err)
	}

	usedMemoryPercent := (float64(memInfo.Used) / float64(memInfo.Total)) * 100.0

	return &models.NodeResources{
		CPUUtilization:    cpuPercent[0],
		MemoryUtilization: usedMemoryPercent,
	}, nil
}
