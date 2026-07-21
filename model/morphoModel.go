package model

type MorphoResponseEntity struct {
	ResponseCode int           `json:"responseCode"`
	ResponseData []VaultEntity `json:"responseData"`
}

type VaultEntity struct {
	Name            string  `json:"name"`
	HealthFactor    float64 `json:"healthFactor"`
	BorrowPnlUsd    float64 `json:"borrowPnlUsd"`
	BorrowAssetsUsd float64 `json:"borrowAssetsUsd"`
	CollateralUsd   float64 `json:"collateralUsd"`
	AvgBorrowApy    float64 `json:"avgBorrowApy"`
	NetBorrowApy    float64 `json:"netBorrowApy"`
	CollateralAsset string  `json:"collateralAsset"`
	LoanAsset       string  `json:"loanAsset"`
}
