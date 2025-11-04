package user

import (
	"swift_transit/domain"
	"swift_transit/rest/handlers/user"
)

type Service interface {
	user.Service //embedding
}

// UserRepo interface
type UserRepo interface {
	Find(userName, password string) (*domain.User, error) // login
	Create(user domain.User) (*domain.User, error)        // create new user
}
