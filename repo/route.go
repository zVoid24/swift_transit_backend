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
	tx, err := r.dbCon.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO routes (name, geom) VALUES ($1, ST_Force2D(ST_GeomFromGeoJSON($2))) RETURNING id`
	var routeID int64
	err = tx.QueryRowx(query, route.Name, route.LineStringGeoJSON).Scan(&routeID)
	if err != nil {
		return nil, err
	}
	route.Id = routeID

	stopQuery := `INSERT INTO stops (route_id, name, stop_order, geom) VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5), 4326)) RETURNING id`
	for i, stop := range route.Stops {
		var stopID int64
		err = tx.QueryRowx(stopQuery, routeID, stop.Name, stop.Order, stop.Lon, stop.Lat).Scan(&stopID)
		if err != nil {
			return nil, err
		}
		route.Stops[i].Id = stopID
		route.Stops[i].RouteId = routeID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &route, nil
}
func (r *routeRepo) FindAll() ([]domain.Route, error) {
	return nil, nil
}
func (r *routeRepo) FindByID(id int64) (*domain.Route, error) {
	return nil, nil
}
