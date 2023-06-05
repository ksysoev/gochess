package gamesrv

import (
	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Game struct {
	ID          string
	Position    string
	PlayerWhite string
	PlayerBlack string
	State       string
}

func NewGame(playerWhite string, playerBlack string) Game {
	id := uuid.New().String()

	newGame := chess.NewGame()
	position := newGame.Position().String()

	return Game{
		ID:          id,
		Position:    position,
		PlayerWhite: playerWhite,
		PlayerBlack: playerBlack,
		State:       "in_progress",
	}
}

func (g *Game) MakeMove(move string) error {
	fen, err := chess.FEN(g.Position)
	if err != nil {
		return err
	}

	game := chess.NewGame(fen)
	err = game.MoveStr(move)
	if err != nil {
		return err
	}

	g.Position = game.Position().String()

	return nil
}
