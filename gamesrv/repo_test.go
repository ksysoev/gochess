package gamesrv

import (
	"testing"
)

func TestGameRepoStorage(t *testing.T) {
	repo := NewGameRepo()

	// Test Add
	game := Game{ID: "1", PlayerWhite: "Alice", PlayerBlack: "Bob"}
	err := repo.Add(game)
	if err != nil {
		t.Errorf("Add returned an error: %v", err)
	}

	// Test Get
	g, err := repo.Get("1")
	if err != nil {
		t.Errorf("Get returned an error: %v", err)
	}
	if g != game {
		t.Errorf("Get returned the wrong game: got %v, want %v", g, game)
	}

	// Test Remove
	err = repo.Remove("1")
	if err != nil {
		t.Errorf("Remove returned an error: %v", err)
	}

	// Test Update
	err = repo.Update(game)
	if err == nil {
		t.Errorf("Update should have returned an error, but it didn't")
	}
}
