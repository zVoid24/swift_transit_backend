package cmd

import (
	"swift_transit/config"
	"swift_transit/infra/db"
	"swift_transit/repo"
	"swift_transit/rest"
	userHandler "swift_transit/rest/handlers/user"
	"swift_transit/rest/middlewares"
	"swift_transit/utils"
	"swift_transit/user"
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
	//repos
	userRepo := repo.NewUserRepo(dbCon)

	//domains
	usrSvc:=user.NewService(userRepo)

	userHandler := userHandler.NewHandler(usrSvc, mngr, utilHandler)
	handler := rest.NewHandler(cnf, middlewareHandler, userHandler)
	handler.Serve()
}
