package measurer

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AndyS1mpson/docker-coscheduler/internal/models"
)

// поиск IPC в выводе команды
var ipcRegexp = `^.*instructions\s*#\s*([\d,]+).*$`

// TaskSpeedMeasurer замеряет "время выполнения" задачи в соответствии с формулой IPC * cpu_time / time
type TaskSpeedMeasurer struct {
	dockerClient dockerClient
}

// NewTaskSpeedMeasurer конструктор для TaskSpeedMeasurer
func NewTaskSpeedMeasurer(dockerClient dockerClient) *TaskSpeedMeasurer {
	return &TaskSpeedMeasurer{
		dockerClient: dockerClient,
	}
}

// MeasureTaskSpeed замеряет время выполнения задачи для FSC
func (t *TaskSpeedMeasurer) Measure(ctx context.Context, task models.Task, duration time.Duration) (time.Duration, error) {
	cpuTime, err := t.dockerClient.GetContainerCPUTime(ctx, task.Config.ContainerID)
	if err != nil {
		return 0, fmt.Errorf("cpu time measure: %w", err)
	}

	ipc, err := measureIPC(duration, task.Config.CPUs.AsString())
	if err != nil {
		return 0, fmt.Errorf("ipc measure: %w", err)
	}

	return ipc * cpuTime, nil
}

// measureIPC измерение IPC с помощью утилиты perf
func measureIPC(duration time.Duration, cpuSet string) (time.Duration, error) {
	args := []string{"stat", "-B"}
	if cpuSet != "" {
		args = append(args, "-C", cpuSet)
	}
	args = append(args, "dd", "if=/dev/zero", "of=/dev/null", fmt.Sprintf("count=%s", duration))

	cmd := exec.Command("perf", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("perf stat: %w", err)
	}

	re := regexp.MustCompile(ipcRegexp)
	match := re.FindStringSubmatch(string(output))
	if len(match) != 2 {
		return 0, fmt.Errorf("can not find IPC in perf stat output")
	}

	ipcStr := strings.ReplaceAll(match[1], ",", ".")
	ipcFloat, err := strconv.ParseFloat(ipcStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ipc: %w", err)
	}

	// Вычисляем IPC в виде временного интервала
	ipc := time.Duration(ipcFloat) * time.Second

	return ipc, nil
}
