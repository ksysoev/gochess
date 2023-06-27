package gamesrv

import (
	"log"

	"github.com/asaskevich/EventBus"
	"github.com/ksysoev/gochess/events"
)

// GameService is a service for managing games.
type GameService struct {
	GameRepo GameRepo
	EventBus EventBus.Bus
}

// NewGameService creates a new GameService.
func NewGameService(gameRepo GameRepo, evbus EventBus.Bus) GameService {
	gs := GameService{
		GameRepo: gameRepo,
		EventBus: evbus,
	}

	evbus.Subscribe("match::found", func(ev events.MatchFoundEvent) {
		// Limitation of the EventBus library
		// we should use go routine here to avoid deadlocks.
		go func() {
			_, err := gs.CreateGame(ev.White, ev.Black)
			if err != nil {
				log.Printf("Error creating game: %s", err)
			}
		}()
	})

	return gs
}

// CreateGame creates a new game.
func (gs GameService) CreateGame(playerWhite string, playerBlack string) (Game, error) {
	game := NewGame(playerWhite, playerBlack)

	err := gs.GameRepo.Add(game)
	if err != nil {
		return Game{}, err
	}

	gs.EventBus.Publish("game:start", events.EventGameStart{
		GameID:      game.ID,
		PlayerWhite: game.PlayerWhite,
		PlayerBlack: game.PlayerBlack,
		Position:    game.Position,
	})

	return game, nil
}

// GetGame gets a game.
func (gs GameService) GetGame(id string) (Game, error) {
	game, err := gs.GameRepo.Get(id)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

// MakeMove makes a move in a game.
func (gs GameService) MakeMove(id string, move string) (Game, error) {
	game, err := gs.GameRepo.Get(id)
	if err != nil {
		return Game{}, err
	}

	err = game.MakeMove(move)
	if err != nil {
		return Game{}, err
	}

	err = gs.GameRepo.Update(game)
	if err != nil {
		return Game{}, err
	}

	gs.EventBus.Publish("game:move", events.EventGameMove{
		GameID:      game.ID,
		Move:        move,
		Position:    game.Position,
		PlayerWhite: game.PlayerWhite,
		PlayerBlack: game.PlayerBlack,
	})
	return game, nil
}
