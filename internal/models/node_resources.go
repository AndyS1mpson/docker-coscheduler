package models

// NodeResources информация о доступности ресурсов ноды
type NodeResources struct {
	CPUUtilization    float64	// Процент загруженности процессора
	MemoryUtilization float64	// Процент загруженности оперативной памяти
}
