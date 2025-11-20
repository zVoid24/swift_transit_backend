package route

import "net/http"

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /route", h.mngr.With(http.HandlerFunc(h.Create)))
}
