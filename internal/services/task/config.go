package task

// NodeConfig конфигурация воркера
type NodeConfig struct {
	CPUNums int64  // Колисество ядер на ноде
	URI     string // URI ноды
	Port    int64  // Порт ноды
}
