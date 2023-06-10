package gamesrv

import (
	"github.com/asaskevich/EventBus"
	"github.com/ksysoev/gochess/events"
)

type GameService struct {
	GameRepo GameRepo
	EventBus EventBus.Bus
}

func NewGameService(gameRepo GameRepo, evbus EventBus.Bus) GameService {
	return GameService{
		GameRepo: gameRepo,
		EventBus: evbus,
	}
}

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

func (gs GameService) GetGame(id string) (Game, error) {
	game, err := gs.GameRepo.Get(id)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

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
