package cmd

import (
	"swift_transit/config"
	"swift_transit/infra/db"
	"swift_transit/repo"
	"swift_transit/rest"
	userHandler "swift_transit/rest/handlers/user"
	"swift_transit/rest/middlewares"
	"swift_transit/user"
	"swift_transit/utils"
)

func Start() {
	cnf := config.Load()
	utilHandler := utils.NewHandler(cnf)
	middlewareHandler := middlewares.NewHandler(utilHandler)
	mngr := middlewareHandler.NewManager()
	dbCon, err := db.NewConnection(&cnf.Db)
	if err != nil {
		panic(err)
	}
	err = db.MigrateDB(dbCon, "./migrations")
	if err != nil {
		panic(err)
	}

	//repos
	userRepo := repo.NewUserRepo(dbCon, utilHandler)

	//domains
	usrSvc := user.NewService(userRepo)

	userHandler := userHandler.NewHandler(usrSvc, middlewareHandler, mngr, utilHandler)
	handler := rest.NewHandler(cnf, middlewareHandler, userHandler)
	handler.Serve()
}
