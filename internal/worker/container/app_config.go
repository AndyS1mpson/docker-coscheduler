package container

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

var configFileName = "worker.yaml"

type NodeConfig struct {
	Host    string `json:"host"`
	Port    int64  `json:"port"`
	CPUNums int64
}

type AppConfig struct {
	WorkerImageHubDir string     `yaml:"image_hub"`
	Node              NodeConfig `yaml:"node"`
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

	config.Node.CPUNums = int64(runtime.NumCPU())

	return config, nil
}
