package gamesrv

type GameService struct {
	GameRepo GameRepo
}

func NewGameService(gameRepo GameRepo) GameService {
	return GameService{
		GameRepo: gameRepo,
	}
}

func (gs GameService) CreateGame(playerWhite string, playerBlack string) (Game, error) {
	game := NewGame(playerWhite, playerBlack)

	err := gs.GameRepo.Add(game)
	if err != nil {
		return Game{}, err
	}

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

	return game, nil
}
