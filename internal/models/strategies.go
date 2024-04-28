package models

var (
	StrategyNameRoundRobin StrategyName = "round-robin"
	StrategyNameFCS        StrategyName = "fcs"
	StrategyNameFCN        StrategyName = "fcn"
)

// StrategyName название стратегии кошедулинга
type StrategyName string

// Strategy описывает стратегию кошедулинга
type Strategy struct {
	ID   int64
	Name StrategyName
}
