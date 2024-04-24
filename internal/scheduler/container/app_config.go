package container

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var configFileName = "scheduler.yaml"

type NodeConfig struct {
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
}

type TaskConfig struct {
	Name       string `yaml:"name"`
	FolderName string `yaml:"folder_name"`
}

// AppConfig структура, содержащая конфигурации менеджеров
type AppConfig struct {
	TaskDir            string        `yaml:"task_dir"`
	TaskInfoDelay      time.Duration `yaml:"task_info_delay"`
	Nodes              []NodeConfig  `yaml:"nodes"`
	Tasks              []TaskConfig  `yaml:"tasks"`
	TaskCombinationNum int64         `yaml:"task_combination_num"`
	FCNTaskNum         int64         `yaml:"fcn_task_num"`
	MeasurementTime    time.Duration `yaml:"measurement_time"`
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
