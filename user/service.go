package user

import (
	"swift_transit/domain"
	"swift_transit/user"
)

type service struct {
	userRepo UserRepo
}

func NewService(usrRepo UserRepo) Service {
	return &service{
		userRepo: usrRepo,
	}
}

func (user.Service) Create(user domain.User) (*domain.User, error)
func (user.Service) Find(username string, password string) (*domain.User, error)
