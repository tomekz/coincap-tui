// fetch data from coincap API
package coincap

import (
	"errors"
	"fmt"
	"time"

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

type AssetHistory struct {
	PriceUsd float64 `json:"priceUsd,string"`
	Time     int64   `json:"time"`
}

type GetAssetHistoryResult struct {
	Data []AssetHistory
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

var apiKey = "2d238e17-b2f0-48b3-a8c9-f3c6fc932c7f" // free public apiKey
var client = resty.New().SetTimeout(10*time.Second).SetHeader("Authorization", "Bearer "+apiKey)

func GetAssets() ([]Asset, error) {
	var result GetAssetsResult
	_, err := client.R().SetResult(&result).Get("https://api.coincap.io/v2/assets")
	if err != nil {
		return nil, DataFetchError{Err: err}
	}

	if len(result.Data) == 0 {
		return nil, DataFetchError{Err: errors.New("no assets found")}
	}
	return result.Data, nil
}

func GetAssetHistory(assetId string) ([]AssetHistory, error) {
	var result GetAssetHistoryResult
	_, err := client.R().SetResult(&result).Get(fmt.Sprintf("https://api.coincap.io/v2/assets/%s/history?interval=d1", assetId))
	if err != nil {
		return nil, DataFetchError{Err: err}
	}

	if len(result.Data) == 0 {
		return nil, DataFetchError{Err: errors.New("no asset history found for assetId " + assetId)}
	}

	// get only most recent history
	chunkSize := len(result.Data) / 2
	res := result.Data[chunkSize:]

	return res, nil
}
