package model

type RiskStatus string

const (
	StatusSafe     RiskStatus = "🟢 Safe"
	StatusMonitor  RiskStatus = "🟡 Monitor"
	StatusCritical RiskStatus = "🔴 Critical"
)

type MetricEvaluation struct {
	Name   string
	Value  string
	Status RiskStatus
	Reason string
}

type RiskReport struct {
	OverallStatus RiskStatus
	Metrics       []MetricEvaluation
	Reason        string
}
