package gamesrv

import "fmt"

type GameRepoStorage struct {
	Storage map[string]Game
}

type GameRepo interface {
	Add(Game) error
	Get(string) (Game, error)
	Remove(string) error
	Update(Game) error
}

func NewGameRepo() GameRepo {
	return GameRepoStorage{
		Storage: make(map[string]Game),
	}
}

func (grs GameRepoStorage) Add(g Game) error {
	if _, ok := grs.Storage[g.ID]; ok {
		return fmt.Errorf("game with id %s already exists", g.ID)
	}

	grs.Storage[g.ID] = g
	return nil
}

func (grs GameRepoStorage) Get(id string) (Game, error) {
	g, ok := grs.Storage[id]
	if !ok {
		return Game{}, fmt.Errorf("game with id %s not found", id)
	}
	return g, nil
}

func (grs GameRepoStorage) Remove(id string) error {
	if _, ok := grs.Storage[id]; !ok {
		return fmt.Errorf("game with id %s not found", id)
	}

	delete(grs.Storage, id)
	return nil
}

func (grs GameRepoStorage) Update(g Game) error {
	if _, ok := grs.Storage[g.ID]; !ok {
		return fmt.Errorf("game with id %s not found", g.ID)
	}

	grs.Storage[g.ID] = g
	return nil
}
