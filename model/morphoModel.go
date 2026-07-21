package model

type MorphoResponseEntity struct {
	ResponseCode int           `json:"responseCode"`
	ResponseData []VaultEntity `json:"responseData"`
}

type VaultEntity struct {
	VaultName     string  `json:"vaultName"`
	TotalAssetUsd float64 `json:"totalAssetUsd"`
	Liquidity     float64 `json:"liquidity"`
	MyAssetUsd    float64 `json:"myAssetUsd"`
	NetApy        float64 `json:"netApy"`
	SharedInVault float64 `json:"sharedInVault"`
}

type MorphoResponseModel struct {
	ResponseCode int          `json:"responseCode"`
	ResponseData []VaultModel `json:"responseData"`
}

type VaultModel struct {
	VaultName     string  `json:"vaultName"`
	TotalAssetUsd float64 `json:"totalAssetUsd"`
	Liquidity     float64 `json:"liquidity"`
	MyAssetUsd    float64 `json:"myAssetUsd"`
	NetApy        float64 `json:"netApy"`
	SharedInVault float64 `json:"sharedInVault"`
}
