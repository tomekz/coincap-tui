// fetch data from coincap API
package coincap

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type Asset struct {
	ID                string  `json:"id"`
	Rank              int64   `json:"rank,string"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	Supply            float64 `json:"supply,string"`
	MaxSupply         float64 `json:"maxSupply,string"`
	MarketCapUsd      float64 `json:"marketCapUsd,string"`
	VolumeUsd24Hr     float64 `json:"volumeUsd24Hr,string"`
	PriceUsd          float64 `json:"priceUsd,string"`
	ChangePercent24Hr float64 `json:"changePercent24Hr,string"`
	Vwap24Hr          float64 `json:"vwap24Hr,string"`
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
