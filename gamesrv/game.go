package gamesrv

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/notnil/chess"
)

// Game represents a chess game.
type Game struct {
	ID          string
	Position    string
	PlayerWhite string
	PlayerBlack string
	State       string
}

// NewGame creates a new game.
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

// MakeMove makes a move in the game.
func (g *Game) MakeMove(move string) error {
	if g.State != "in_progress" {
		return fmt.Errorf("Game is finished")
	}

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

	switch game.Outcome() {
	case chess.WhiteWon:
		g.State = "white_won"
	case chess.BlackWon:
		g.State = "black_won"
	case chess.Draw:
		g.State = "draw"
	default:
		g.State = "in_progress"
	}

	return nil
}
