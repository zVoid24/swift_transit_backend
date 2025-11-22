package domain

type Route struct {
	Id                int64  `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	LineStringGeoJSON string `json:"linestring_geojson" db:"-"` // Not stored directly, used for geom insertion
	Stops             []Stop `json:"stops"`
}
