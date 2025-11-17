package cmd

import (
	"context"
	"fmt"
	"swift_transit/config"
	"swift_transit/infra/db"
	redisConf "swift_transit/infra/redis"
	"swift_transit/repo"
	"swift_transit/rest"
	userHandler "swift_transit/rest/handlers/user"
	"swift_transit/rest/middlewares"
	"swift_transit/user"
	"swift_transit/utils"
)

func Start() {
	ctx := context.Background()
	cnf := config.Load()
	utilHandler := utils.NewHandler(cnf)
	middlewareHandler := middlewares.NewHandler(utilHandler)
	mngr := middlewareHandler.NewManager()
	redisCon, err := redisConf.NewConnection(&cnf.RedisCnf, ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(redisCon)
	// err = redisCon.Set(ctx, "name", "zahid", 0).Err()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }
	// val, err := redisCon.Get(ctx, "name").Result()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	panic(err)
	// }
	// fmt.Println(val)
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
