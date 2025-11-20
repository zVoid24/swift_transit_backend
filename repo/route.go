package repo

import (
	"swift_transit/route"
	"swift_transit/utils"

	"github.com/jmoiron/sqlx"
)

type RouteRepo interface {
	route.RouteRepo
}

type routeRepo struct {
	dbCon       *sqlx.DB
	utilHandler *utils.Handler
}

func NewRouteRepo(dbcon *sqlx.DB, utilHandler *utils.Handler) RouteRepo {
	return &routeRepo{
		dbCon:       dbcon,
		utilHandler: utilHandler,
	}
}
