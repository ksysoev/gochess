package users

import "fmt"

type UserRepoStorage struct {
	Storage map[string]User
}

// GameRepo is an interface for a game repository.
type UserRepo interface {
	Add(User) error
	Get(string) (User, error)
	Remove(string) error
	Update(User) error
}

// NewGameRepo returns a new UserRepoStorage
func NewUserRepo() UserRepo {
	return UserRepoStorage{
		Storage: make(map[string]User),
	}
}

// Add adds a user to the repository.
func (grs UserRepoStorage) Add(g User) error {
	if _, ok := grs.Storage[g.ID]; ok {
		return fmt.Errorf("user with id %s already exists", g.ID)
	}

	grs.Storage[g.ID] = g
	return nil
}

// Get returns a user from the repository.
func (grs UserRepoStorage) Get(id string) (User, error) {
	g, ok := grs.Storage[id]
	if !ok {
		return User{}, fmt.Errorf("user with id %s not found", id)
	}
	return g, nil
}

// Remove removes a user from the repository.
func (grs UserRepoStorage) Remove(id string) error {
	if _, ok := grs.Storage[id]; !ok {
		return fmt.Errorf("user with id %s not found", id)
	}

	delete(grs.Storage, id)
	return nil
}

// Update updates a user in the repository.
func (grs UserRepoStorage) Update(g User) error {
	if _, ok := grs.Storage[g.ID]; !ok {
		return fmt.Errorf("user with id %s not found", g.ID)
	}

	grs.Storage[g.ID] = g
	return nil
}
