package rest

import (
	"swift_transit/config"
	"swift_transit/rest/handlers/user"
	"swift_transit/rest/middlewares"
)

type Handler struct {
	cnf         *config.Config
	mdlw        *middlewares.Handler
	userHandler *user.Handler
}

func NewHandler(cnf *config.Config, mdlw *middlewares.Handler, userHandler *user.Handler) *Handler {
	return &Handler{
		cnf:         cnf,
		mdlw:        mdlw,
		userHandler: userHandler,
	}
}
