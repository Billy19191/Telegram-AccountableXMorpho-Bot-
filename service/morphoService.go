package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
)

type MorphoService struct {
	BaseURL       string
	WalletAddress string
	ChainID       string
	HttpClient    *http.Client
}

func NewMorphoService(baseURL, walletAddress, chainID string) *MorphoService {
	return &MorphoService{
		BaseURL:       baseURL,
		WalletAddress: walletAddress,
		ChainID:       chainID,
		HttpClient:    &http.Client{},
	}
}

func (s *MorphoService) GetBorrowPosition() (*model.MorphoResponseEntity, error) {
	url := fmt.Sprintf("%s/api/v1/borrowPositionPosition?walletAddress=%s&chainID=%s", s.BaseURL, s.WalletAddress, s.ChainID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vault positions: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("failed to fetch vault positions: unexpected status code %d", response.StatusCode)
	}
	var result model.MorphoResponseEntity
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
