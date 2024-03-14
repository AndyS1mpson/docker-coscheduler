package models

// Node описывает ноду, на которой выполняются таски
type Node struct {
	Port    int64  `json:"port"`     // Порт ноды
	Host    string `json:"host"`     // Хост ноды
	CPUNums int64  `json:"cpu_nums"` // Количество доступных ядер
}
