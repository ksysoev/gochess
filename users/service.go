package users

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo UserRepo
}

func NewUserService(ur UserRepo) UserService {
	return UserService{
		UserRepo: ur,
	}
}

func (us UserService) SignUp(username string, password string) (*User, error) {
	id := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := User{
		ID:           id,
		UserName:     username,
		PasswordHash: string(hashedPassword),
	}

	err = us.UserRepo.Add(u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (us UserService) SignIn(username string, password string) (*User, error) {
	//TODO: here we need to get user by username... but this methods use id
	u, err := us.UserRepo.Get(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return &u, nil
}
