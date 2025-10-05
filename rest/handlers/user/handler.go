package user

import "swift_transit/repo"

type Handler struct {
	UserRepo repo.UserRepo
}

func NewHandler(userRepo repo.UserRepo) *Handler {
	return &Handler{
		UserRepo: userRepo,
	}
}
