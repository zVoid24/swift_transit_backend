package user

import (
	"swift_transit/repo"
	"swift_transit/rest/middlewares"
	"swift_transit/utils"
)

type Handler struct {
	UserRepo    repo.UserRepo
	mngr        *middlewares.Manager
	utilHandler *utils.Handler
}

func NewHandler(userRepo repo.UserRepo, mngr *middlewares.Manager, utilHandler *utils.Handler) *Handler {
	return &Handler{
		UserRepo:    userRepo,
		mngr:        mngr,
		utilHandler: utilHandler,
	}
}
