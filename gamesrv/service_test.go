package gamesrv

import (
	"errors"
	"testing"

	"github.com/asaskevich/EventBus"
)

type mockGameRepo struct {
	games map[string]Game
}

func (m *mockGameRepo) Add(game Game) error {
	if _, ok := m.games[game.ID]; ok {
		return errors.New("game already exists")
	}
	m.games[game.ID] = game
	return nil
}

func (m *mockGameRepo) Get(id string) (Game, error) {
	if game, ok := m.games[id]; ok {
		return game, nil
	}
	return Game{}, errors.New("game not found")
}

func (m *mockGameRepo) Update(game Game) error {
	if _, ok := m.games[game.ID]; !ok {
		return errors.New("game not found")
	}
	m.games[game.ID] = game
	return nil
}

func (m *mockGameRepo) Remove(id string) error {
	if _, ok := m.games[id]; !ok {
		return errors.New("game not found")
	}
	delete(m.games, id)
	return nil
}

func TestGameServiceCreateGame(t *testing.T) {
	mockRepo := &mockGameRepo{games: make(map[string]Game)}
	service := NewGameService(mockRepo, EventBus.New())

	playerWhite := "Alice"
	playerBlack := "Bob"

	game, err := service.CreateGame(playerWhite, playerBlack)
	if err != nil {
		t.Errorf("CreateGame returned an error: %v", err)
	}

	if game.PlayerWhite != playerWhite {
		t.Errorf("CreateGame returned a game with the wrong playerWhite: got %v, want %v", game.PlayerWhite, playerWhite)
	}

	if game.PlayerBlack != playerBlack {
		t.Errorf("CreateGame returned a game with the wrong playerBlack: got %v, want %v", game.PlayerBlack, playerBlack)
	}

	if _, ok := mockRepo.games[game.ID]; !ok {
		t.Errorf("CreateGame didn't add the game to the repository")
	}
}

func TestGameServiceGetGame(t *testing.T) {
	mockRepo := &mockGameRepo{games: make(map[string]Game)}
	service := NewGameService(mockRepo, EventBus.New())

	playerWhite := "Alice"
	playerBlack := "Bob"

	game, _ := service.CreateGame(playerWhite, playerBlack)

	g, err := service.GetGame(game.ID)
	if err != nil {
		t.Errorf("GetGame returned an error: %v", err)
	}

	if g != game {
		t.Errorf("GetGame returned the wrong game: got %v, want %v", g, game)
	}
}

func TestGameServiceMakeMove(t *testing.T) {
	mockRepo := &mockGameRepo{games: make(map[string]Game)}
	service := NewGameService(mockRepo, EventBus.New())

	playerWhite := "Alice"
	playerBlack := "Bob"

	game, _ := service.CreateGame(playerWhite, playerBlack)

	_, err := service.MakeMove(game.ID, "e2e4")
	if err != nil {
		t.Errorf("MakeMove returned an error: %v", err)
	}

	g, err := service.GetGame(game.ID)
	if err != nil {
		t.Errorf("GetGame returned an error: %v", err)
	}

	if g.Position != "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1" {
		t.Errorf("MakeMove didn't update the game position correctly: got %v, want %v", g.Position, "rnbqkbnr/pppp1ppp/4p3/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1")
	}

	_, err = service.MakeMove(game.ID, "e7e5")
	if err != nil {
		t.Errorf("MakeMove returned an error: %v", err)
	}

	g, err = service.GetGame(game.ID)
	if err != nil {
		t.Errorf("GetGame returned an error: %v", err)
	}

	if g.Position != "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2" {
		t.Errorf("MakeMove didn't update the game position correctly: got %v, want %v", g.Position, "rnbqkbnr/ppp2ppp/3pp3/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2")
	}

	_, err = service.MakeMove(game.ID, "g1f3")
	if err != nil {
		t.Errorf("MakeMove returned an error: %v", err)
	}

	g, err = service.GetGame(game.ID)
	if err != nil {
		t.Errorf("GetGame returned an error: %v", err)
	}

	if g.Position != "rnbqkbnr/pppp1ppp/8/4p3/4P3/5P2/PPPP2PP/RNBQKBNR b KQkq - 0 2" {
		t.Errorf("MakeMove didn't update the game position correctly: got %v, want %v", g.Position, "rnbqkbnr/ppp2ppp/3pp3/8/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2")
	}

	_, err = service.MakeMove(game.ID, "d8h4")
	if err == nil {
		t.Errorf("MakeMove should have returned an error, but it didn't")
	}
}
