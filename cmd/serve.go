package cmd

import (
	"swift_transit/config"
	"swift_transit/infra/db"
	"swift_transit/repo"
	"swift_transit/rest"
	"swift_transit/rest/handlers/user"
	"swift_transit/rest/middlewares"
	"swift_transit/utils"
)

func Start() {
	cnf := config.Load()
	middlewareHandler := middlewares.NewHandler()
	mngr := middlewareHandler.NewManager()
	dbCon, err := db.NewConnection(&cnf.Db)
	if err != nil {
		panic(err)
	}
	err = db.MigrateDB(dbCon, "./migrations")
	if err != nil {
		panic(err)
	}
	utilHandler := utils.NewHandler(cnf)
	userRepo := repo.NewUserRepo(dbCon)
	userHandler := user.NewHandler(userRepo, mngr, utilHandler)
	handler := rest.NewHandler(cnf, middlewareHandler, userHandler)
	handler.Serve()
}
