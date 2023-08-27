// fetch data from coincap API
package coincap

import (
	"errors"
	"fmt"
	"time"

	resty "github.com/go-resty/resty/v2"
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

var client = resty.New().SetTimeout(5 * time.Second)

func GetAssets(limit int) ([]Asset, error) {
	log.Debugf("Fetching assets from coincap API")
	var result GetAssetsResult
	_, err := client.R().SetResult(&result).Get(fmt.Sprintf("https://api.coincap.io/v2/assets?limit=%d", limit))
	if err != nil {
		return nil, DataFetchError{Err: err}
	}

	if len(result.Data) == 0 {
		return nil, DataFetchError{Err: errors.New("no assets found")}
	}
	return result.Data, nil
}

func GetAssetHistory(assetId string) ([]AssetHistory, error) {
	log.Debugf("Fetching asset history from coincap API for assetId %s", assetId)
	var result GetAssetHistoryResult
	_, err := client.R().SetResult(&result).Get(fmt.Sprintf("https://api.coincap.io/v2/assets/%s/history?interval=h6", assetId))
	if err != nil {
		return nil, DataFetchError{Err: err}
	}

	if len(result.Data) == 0 {
		return nil, DataFetchError{Err: errors.New("no asset history found for assetId " + assetId)}
	}

	// workaround to get last 14 days of data. `start` and `end` API params are not working
	chunkSize := len(result.Data) / 7
	res := result.Data[6*chunkSize:]

	return res, nil
}
