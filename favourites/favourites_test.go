package favourites

import (
	"testing"
)

func TestNew(t *testing.T) {
	f := New()
	if f == nil {
		t.Error("New returned nil")
	}
	if len(f.favs) != 0 {
		t.Error("New did not create an empty favourites map")
	}
}

func TestFavourite(t *testing.T) {
	f := New()

	// Add asset to favourites
	f.Favourite("abc")
	if !f.favs["abc"] {
		t.Error("Favourite did not add asset to favourites")
	}

	// Remove asset from favourites
	f.Favourite("abc")
	if f.favs["abc"] {
		t.Error("Favourite did not remove asset from favourites")
	}

	// Add multiple assets to favourites
	f.Favourite("def")
	f.Favourite("ghi")
	if !f.favs["def"] || !f.favs["ghi"] {
		t.Error("Favourite did not add multiple assets to favourites")
	}
}
