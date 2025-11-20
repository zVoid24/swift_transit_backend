package repo

import (
	"swift_transit/domain"
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

func (r *routeRepo) Create(route domain.Route) (*domain.Route, error) {
	return nil, nil
}
func (r *routeRepo) FindAll() ([]domain.Route, error) {
	return nil, nil
}
func (r *routeRepo) FindByID(id int64) (*domain.Route, error) {
	return nil, nil
}
