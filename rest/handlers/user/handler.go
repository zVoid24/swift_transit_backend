package user

import (
	"swift_transit/rest/middlewares"
	"swift_transit/utils"
)

type Handler struct {
	svc         Service
	mngr        *middlewares.Manager
	utilHandler *utils.Handler
}

func NewHandler(svc Service, mngr *middlewares.Manager, utilHandler *utils.Handler) *Handler {
	return &Handler{
		svc:         svc,
		mngr:        mngr,
		utilHandler: utilHandler,
	}
}
