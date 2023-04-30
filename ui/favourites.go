package ui

import (
	"encoding/json"
	"os"
	"path/filepath"

	gap "github.com/muesli/go-app-paths"
)

const FILE_NAME = "fav_assets.json"

/* Map of favourites */
var Favs map[string]bool = make(map[string]bool)

/* Adds or removes asset from favourites */
func favourite(assetId string) {
	exists := Favs[assetId]
	if exists {
		delete(Favs, assetId)
	} else {
		Favs[assetId] = true
	}
}

/*
Saves favourites to JSON file. Creates file if it doesn't exist.
Uses the plaform specific aplication data directory.
*/
func saveFavourites() error {
	path, err := getFilePath(FILE_NAME)
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o754); err != nil {
			return err
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		content, err := json.Marshal(Favs)
		if err != nil {
			return err
		}
		os.WriteFile(path, content, 0o644)
	}
	return nil
}

/* loads favourites from file */
func loadFavourites() error {
	path, err := getFilePath(FILE_NAME)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(content, &Favs)
	if err != nil {
		return err
	}
	return nil
}

/* returns the full path to a file in the application's default data directory */
func getFilePath(name string) (string, error) {
	home := gap.NewScope(gap.User, "coincap-tui")
	path, err := home.DataPath(name)
	if err != nil {
		return "", err
	}

	return path, err
}
