package gamesrv

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	playerWhite := "Alice"
	playerBlack := "Bob"

	game := NewGame(playerWhite, playerBlack)

	if game.PlayerWhite != playerWhite {
		t.Errorf("NewGame returned a game with the wrong playerWhite: got %v, want %v", game.PlayerWhite, playerWhite)
	}

	if game.PlayerBlack != playerBlack {
		t.Errorf("NewGame returned a game with the wrong playerBlack: got %v, want %v", game.PlayerBlack, playerBlack)
	}

	if game.State != "in_progress" {
		t.Errorf("NewGame returned a game with the wrong state: got %v, want %v", game.State, "in_progress")
	}
}

func TestMakeMove(t *testing.T) {
	playerWhite := "Alice"
	playerBlack := "Bob"

	game := NewGame(playerWhite, playerBlack)

	err := game.MakeMove("e4")
	if err != nil {
		t.Errorf("MakeMove returned an error: %v", err)
	}

	if game.Position != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Errorf("MakeMove didn't update the game position correctly: got %v, want %v", game.Position, "rnbqkbnr/pppp1ppp/4p3/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1")
	}

	err = game.MakeMove("e5")
	if err != nil {
		t.Errorf("MakeMove returned an error: %v", err)
	}

	if game.Position != "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2" {
		t.Errorf("MakeMove didn't update the game position correctly: got %v, want %v", game.Position, "rnbqkbnr/ppp2ppp/3pp3/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2")
	}

	err = game.MakeMove("g1f3")
	if err != nil {
		t.Errorf("MakeMove returned an error: %v", err)
	}

	if game.Position != "rnbqkbnr/pppp1ppp/8/4p3/4P3/5P2/PPPP2PP/RNBQKBNR b KQkq - 0 2" {
		t.Errorf("MakeMove didn't update the game position correctly: got %v, want %v", game.Position, "rnbqkbnr/ppp2ppp/3pp3/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2")
	}

	err = game.MakeMove("d8h4")
	if err == nil {
		t.Errorf("MakeMove should have returned an error, but it didn't")
	}
}
