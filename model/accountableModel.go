package model

type AccountableResponseEntity struct {
	ResponseCode int             `json:"responseCode"`
	ResponseData ValueAllocation `json:"responseData"`
}

type ValueAllocation struct {
	VaultAllocations []AccountableVaultAllocationEntity `json:"vault_allocations"`
}
type AccountableVaultAllocationEntity struct {
	Apy                      float64 `json:"apy"`
	CanDeposit               bool    `json:"canDeposit"`
	ChainID                  int     `json:"chainId"`
	Company                  string  `json:"company"`
	CompanyWordmarkLogo      string  `json:"company_wordmark_logo"`
	Currency                 string  `json:"currency"`
	Duration                 int     `json:"duration"`
	EstablishmentFeeToPayUsd float64 `json:"establishmentFeeToPayUsd"`
	Implementation           string  `json:"implementation"`
	InterestPaymentToPayUsd  float64 `json:"interestPaymentToPayUsd"`
	Label                    string  `json:"label"`
	LifeCurrentDays          int     `json:"lifeCurrentDays"`
	LoanAddress              string  `json:"loanAddress"`
	LoanID                   int     `json:"loanId"`
	MyDepositUsd             float64 `json:"myDepositUsd"`
	PermissionLevel          int     `json:"permissionLevel"`
	Pnl                      float64 `json:"pnl"`
	Schedule                 []any   `json:"schedule"`
	TotalInterestClaimedUsd  float64 `json:"totalInterestClaimedUsd"`
	TotalToPayUsd            float64 `json:"totalToPayUsd"`
	UnrealizedPnl            float64 `json:"unrealizedPnl"`
	Value                    float64 `json:"value"`
	VaultAddress             string  `json:"vaultAddress"`
	VaultName                string  `json:"vaultName"`
}
