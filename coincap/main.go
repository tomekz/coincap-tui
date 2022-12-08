// fetch data from coincap API
package coincap

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type Asset struct {
	ID                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
}

type GetAssetsResult struct {
	Data []Asset
}

type DataFetchError struct {
	Err error
}

func (e DataFetchError) Error() string {
	return e.Err.Error()
}

func GetAssets() ([]Asset, error) {
	var result GetAssetsResult
	client := resty.New()
	_, err := client.R().SetResult(&result).Get("https://api.coincap.io/v2/assets")
	if err != nil {
		return nil, DataFetchError{Err: err}
	}

	if len(result.Data) == 0 {
		return nil, DataFetchError{Err: errors.New("no assets found")}
	}
	return result.Data, nil
}
