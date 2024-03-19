package models

// StrategyTask описывает задачу, поступающую планировщику на вход
type StrategyTask struct {
	Name       string // Название задачи
	FolderName string // Название папки, в которой лежат все файлы задачи
}
