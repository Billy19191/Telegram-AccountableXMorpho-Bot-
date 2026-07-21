package service

import (
	"fmt"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
)

func EvaluateVaultRisk(accountable model.AccountableVaultAllocationEntity, vault model.VaultEntity) model.RiskReport {
	metrics := []model.MetricEvaluation{
		evaluateLiquidityRatio(vault),
		evaluateShareConcentration(vault),
		evaluateAPY(accountable.Apy, vault.NetApy),
		evaluateVaultUtilization(vault),
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

func evaluateLiquidityRatio(vault model.VaultEntity) model.MetricEvaluation {
	var ratio float64
	if vault.MyAssetUsd > 0 {
		ratio = vault.Liquidity / vault.MyAssetUsd
	}

	var status model.RiskStatus
	var reason string

	switch {
	case ratio < 3:
		status = model.StatusCritical
		reason = "Liquidity is less than 3x your position — withdrawal may be difficult"
	case ratio < 10:
		status = model.StatusMonitor
		reason = "Liquidity is below 10x your position — keep monitoring"
	default:
		status = model.StatusSafe
		reason = "Sufficient liquidity available for your position"
	}

	return model.MetricEvaluation{
		Name:   "Liquidity Ratio",
		Value:  fmt.Sprintf("%.1fx", ratio),
		Status: status,
		Reason: reason,
	}
}

func evaluateShareConcentration(vault model.VaultEntity) model.MetricEvaluation {
	share := vault.SharedInVault

	var status model.RiskStatus
	var reason string

	switch {
	case share > 5:
		status = model.StatusCritical
		reason = "You hold over 5% of the vault — high concentration risk"
	case share > 2:
		status = model.StatusMonitor
		reason = "You hold over 2% of the vault — moderate concentration"
	default:
		status = model.StatusSafe
		reason = "Your share in the vault is well diversified"
	}

	return model.MetricEvaluation{
		Name:   "Share in Vault",
		Value:  fmt.Sprintf("%.2f%%", share),
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
	case ratio >= 70:
		status = model.StatusCritical
		reason = "Borrow APY is at or above 70% of deposit APY"
	case ratio >= 50:
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

func evaluateVaultUtilization(vault model.VaultEntity) model.MetricEvaluation {
	var utilization float64
	if vault.TotalAssetUsd > 0 {
		utilization = ((vault.TotalAssetUsd - vault.Liquidity) / vault.TotalAssetUsd) * 100
	}

	var status model.RiskStatus
	var reason string

	switch {
	case utilization > 90:
		status = model.StatusCritical
		reason = "Over 90% of vault assets are utilized — very low available liquidity"
	case utilization > 80:
		status = model.StatusMonitor
		reason = "Utilization exceeds 80% — liquidity may tighten"
	default:
		status = model.StatusSafe
		reason = "Vault utilization is at a healthy level"
	}

	return model.MetricEvaluation{
		Name:   "Vault Utilization",
		Value:  fmt.Sprintf("%.2f", utilization),
		Status: status,
		Reason: reason,
	}
}
