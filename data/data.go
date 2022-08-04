package data

import (
	"errors"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Asset struct {
	AssetID            string    `json:"asset_id"`
	Name               string    `json:"name"`
	TypeIsCrypto       int       `json:"type_is_crypto"`
	DataStart          string    `json:"data_start"`
	DataEnd            string    `json:"data_end"`
	DataQuoteStart     time.Time `json:"data_quote_start"`
	DataQuoteEnd       time.Time `json:"data_quote_end"`
	DataOrderbookStart time.Time `json:"data_orderbook_start"`
	DataOrderbookEnd   time.Time `json:"data_orderbook_end"`
	DataTradeStart     time.Time `json:"data_trade_start"`
	DataTradeEnd       time.Time `json:"data_trade_end"`
	DataSymbolsCount   int       `json:"data_symbols_count"`
	Volume1HrsUsd      float64   `json:"volume_1hrs_usd"`
	Volume1DayUsd      float64   `json:"volume_1day_usd"`
	Volume1MthUsd      float64   `json:"volume_1mth_usd"`
	PriceUsd           float64   `json:"price_usd"`
}
type DataFetchError struct {
	Err error
}

func (e DataFetchError) Error() string {
	return e.Err.Error()
}

var apiKey = os.Getenv("API_KEY")

func SearchAssets(asset string) ([]Asset, error) {
	var assets []Asset
	client := resty.New()
	_, err := client.R().SetHeader("X-CoinAPI-Key", apiKey).SetResult(&assets).Get("https://rest.coinapi.io/v1/assets/" + asset)
	if err != nil {
		return nil, DataFetchError{Err: err}
	}

	if len(assets) == 0 {
		return nil, DataFetchError{Err: errors.New("no assets found")}
	}
	return assets, nil
}
