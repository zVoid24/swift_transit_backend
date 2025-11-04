package user

import "swift_transit/domain"

type Service interface {
	Find(username string, password string) (*domain.User, error)
	Create(user domain.User) (*domain.User, error)
}
