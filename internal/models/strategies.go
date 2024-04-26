package models

var (
	StrategyNameRoundRobin StrategyName = "round-robin"
	StrategyNameFCS        StrategyName = "fcs"
	StrategyNameFCN        StrategyName = "fcn"
)

type StrategyName string
