package repo

import (
	"fmt"
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
	var route domain.Route
	query := `SELECT id, name, ST_AsGeoJSON(geom) as linestring_geojson FROM routes WHERE id = $1`
	err := r.dbCon.Get(&route, query, id)
	if err != nil {
		return nil, err
	}

	var stops []domain.Stop
	stopQuery := `SELECT id, route_id,stop_order,name, ST_X(geom::geometry) as lon, ST_Y(geom::geometry) as lat FROM stops WHERE route_id = $1 ORDER BY stop_order`

	err = r.dbCon.Select(&stops, stopQuery, id)
	fmt.Println(stops)
	if err != nil {
		return nil, err
	}

	route.Stops = stops

	return &route, nil
}
