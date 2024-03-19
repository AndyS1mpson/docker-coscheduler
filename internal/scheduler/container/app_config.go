package container

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var configFileName = "scheduler.yaml"

type NodeConfig struct {
	Host string `json:"host"`
	Port int64  `json:"port"`
}

// AppConfig структура, содержащая конфигурации менеджеров
type AppConfig struct {
	TaskDir       string        `yaml:"task_dir"`
	TaskInfoDelay time.Duration `yaml:"task_info_delay"`
	Nodes         []NodeConfig  `yaml:"nodes"`
}

// NewConfig returns a new decoded Config struct
func NewConfig() (*AppConfig, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(rootDir, configFileName)

	config := &AppConfig{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
