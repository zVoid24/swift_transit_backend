package rest

import (
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) Serve() {
	mux := http.NewServeMux()
	mngr := h.mdlw.NewManager()
	mngr.Use(h.mdlw.Logger, h.mdlw.Cors)
	wrappedMux := mngr.WrapMux(mux)
	//InitRoutes(mux, *mngr)
	fmt.Println("Server running on", h.cnf.HttpPort)
	port := ":" + h.cnf.HttpPort
	err := http.ListenAndServe(port, wrappedMux)
	if err != nil {
		log.Fatal("Error starting server")
	}
}
