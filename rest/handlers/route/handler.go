package route

import "swift_transit/rest/middlewares"

type Handler struct {
	mngr middlewares.Manager
}

func NewHandler(mngr middlewares.Manager) *Handler {
	return &Handler{
		mngr: mngr,
	}
}
