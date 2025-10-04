package cmd

import (
	"swift_transit/config"
	"swift_transit/rest"
	"swift_transit/rest/middlewares"
)

func Start() {
	cnf := config.Load()
	middlewareHandler := middlewares.NewHandler()
	handler := rest.NewHandler(cnf, middlewareHandler)
	handler.Serve()
}
