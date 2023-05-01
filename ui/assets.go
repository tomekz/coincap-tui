package ui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/thoas/go-funk"
	"github.com/tomekz/coincap-tui/coincap"
)

type Asset struct {
	ID                string
	Rank              int64
	Symbol            string
	Name              string
	Supply            float64
	MaxSupply         float64
	MarketCapUsd      float64
	VolumeUsd24Hr     float64
	PriceUsd          float64
	ChangePercent24Hr float64
	Vwap24Hr          float64
	IsFavourite       bool
}

func mapAssets(assets []coincap.Asset) []Asset {
	return funk.Map(assets, func(asset coincap.Asset) Asset {
		return Asset{
			ID:                asset.ID,
			Rank:              asset.Rank,
			Symbol:            asset.Symbol,
			Name:              asset.Name,
			Supply:            asset.Supply,
			MaxSupply:         asset.MaxSupply,
			MarketCapUsd:      asset.MarketCapUsd,
			VolumeUsd24Hr:     asset.VolumeUsd24Hr,
			PriceUsd:          asset.PriceUsd,
			ChangePercent24Hr: asset.ChangePercent24Hr,
			Vwap24Hr:          asset.Vwap24Hr,
			IsFavourite:       Favs.Get()[asset.ID],
		}
	}).([]Asset)
}

func updateFavourite(assets []Asset, assetId string) []Asset {
	return funk.Map(assets, func(asset Asset) Asset {
		if asset.ID == assetId {
			asset.IsFavourite = !asset.IsFavourite
		}
		return asset
	}).([]Asset)
}

func getFavouriteAssets(assets []Asset) []Asset {
	favourites := funk.Filter(assets, func(asset Asset) bool {
		return asset.IsFavourite
	})
	return favourites.([]Asset)
}

func getRows(assets []Asset) []table.Row {
	return funk.Map(assets, func(asset Asset) table.Row {
		fav := ""
		if Favs.Get()[asset.ID] {
			fav = "★"
		}
		maxSupply := formatFloat(asset.MaxSupply)

		if asset.MaxSupply == 0 {
			maxSupply = "∞"
		}
		return []string{
			strconv.FormatInt(asset.Rank, 10),
			fav,
			asset.ID,
			asset.Symbol,
			strconv.FormatFloat(asset.PriceUsd, 'f', 2, 64),
			formatPercent(asset.ChangePercent24Hr),
			formatFloat(asset.Supply),
			maxSupply,
			formatFloat(asset.MarketCapUsd),
			formatFloat(asset.VolumeUsd24Hr),
		}
	}).([]table.Row)
}
