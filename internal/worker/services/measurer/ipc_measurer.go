package measurer

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// поиск IPC в выводе команды
var ipcRegexp = `^.*instructions\s*#\s*([\d,]+).*$`

// MeasureIPC измерение IPC с помощью утилиты perf
func MeasureIPC(duration string, cpuSet string) (float64, error) {
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
	ipc, err := strconv.ParseFloat(ipcStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parse ipc: %w", err)
	}

	return ipc, nil
}
