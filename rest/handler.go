package rest

import (
	"swift_transit/config"
	"swift_transit/rest/middlewares"
)

type Handler struct {
	cnf  *config.Config
	mdlw *middlewares.Handler
}

func NewHandler(cnf *config.Config, mdlw *middlewares.Handler) *Handler {
	return &Handler{
		cnf:  cnf,
		mdlw: mdlw,
	}
}
