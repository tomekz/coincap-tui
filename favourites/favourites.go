package favourites

import (
	"encoding/json"
	"os"
	"path/filepath"

	gap "github.com/muesli/go-app-paths"
)

const FILE_NAME = "fav_assets.json"

type Favourites struct {
	// map of fovourites
	favs map[string]bool
}

// New returns a new Favourites struct
func New() *Favourites {
	return &Favourites{
		favs: make(map[string]bool),
	}
}

// Adds or removes asset from favourites
func (f *Favourites) Favourite(assetId string) {
	exists := f.favs[assetId]
	if exists {
		delete(f.favs, assetId)
	} else {
		f.favs[assetId] = true
	}
}

// Saves favourites to JSON file. Creates file if it doesn't exist.
//
// Uses the plaform specific aplication data directory.
func (f *Favourites) Save() error {
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
		content, err := json.Marshal(f.favs)
		if err != nil {
			return err
		}
		//nolint:errcheck
		os.WriteFile(path, content, 0o644)
	}
	return nil
}

// Loads favourites from file
func (f *Favourites) Load() error {
	path, err := getFilePath(FILE_NAME)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(content, &f.favs)
	if err != nil {
		return err
	}
	return nil
}

func (f *Favourites) Get() map[string]bool {
	return f.favs
}

// returns the full path to a file in the application's default data directory
func getFilePath(name string) (string, error) {
	home := gap.NewScope(gap.User, "coincap-tui")
	path, err := home.DataPath(name)
	if err != nil {
		return "", err
	}

	return path, err
}
