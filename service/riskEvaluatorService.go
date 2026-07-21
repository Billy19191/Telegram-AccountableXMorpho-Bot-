package service

import (
	"fmt"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
)

func EvaluateVaultRisk(accountable model.AccountableVaultAllocationEntity, vault model.VaultEntity) model.RiskReport {
	metrics := []model.MetricEvaluation{
		evaluateHealthFactor(vault),
		evaluateAPY(accountable.Apy, vault.NetBorrowApy),
	}

	overall := model.StatusSafe
	reason := ""

	for _, m := range metrics {
		if m.Status == model.StatusCritical {
			overall = model.StatusCritical
			if reason != "" {
				reason += "\n"
			}
			reason += m.Reason
			continue
		}
		if m.Status == model.StatusMonitor {
			if overall != model.StatusCritical {
				overall = model.StatusMonitor
			}
			if reason != "" {
				reason += "\n"
			}
			reason += m.Reason
		}
	}

	return model.RiskReport{
		OverallStatus: overall,
		Metrics:       metrics,
		Reason:        reason,
	}
}

func evaluateHealthFactor(vault model.VaultEntity) model.MetricEvaluation {
	var status model.RiskStatus
	var reason string

	switch {
	case vault.HealthFactor <= 0:
		status = model.StatusMonitor
		reason = "Health factor is unavailable"
	case vault.HealthFactor < 1.1:
		status = model.StatusCritical
		reason = "Health factor is below 1.1 — position is at elevated risk"
	case vault.HealthFactor < 1.5:
		status = model.StatusMonitor
		reason = "Health factor is below 1.5 — keep monitoring"
	default:
		status = model.StatusSafe
		reason = "Health factor is healthy"
	}

	return model.MetricEvaluation{
		Name:   "Health Factor",
		Value:  fmt.Sprintf("%.2f", vault.HealthFactor),
		Status: status,
		Reason: reason,
	}
}

func evaluateAPY(depositApy, borrowApy float64) model.MetricEvaluation {
	ratio := 0.0
	if depositApy > 0 {
		ratio = (borrowApy / depositApy) * 100
	}

	var status model.RiskStatus
	var reason string

	switch {
	case depositApy <= 0:
		status = model.StatusMonitor
		reason = "Deposit APY is unavailable, so borrow-to-deposit risk could not be assessed"
	case ratio >= 80:
		status = model.StatusCritical
		reason = "Borrow APY is at or above 80% of deposit APY"
	case ratio >= 70:
		status = model.StatusMonitor
		reason = "Borrow APY is approaching 70% of deposit APY"
	default:
		status = model.StatusSafe
		reason = "Borrow APY remains well below deposit APY"
	}

	return model.MetricEvaluation{
		Name:   "Borrow vs Deposit APY",
		Value:  fmt.Sprintf("%.2f%% / %.2f%%", borrowApy, depositApy),
		Status: status,
		Reason: reason,
	}
}
