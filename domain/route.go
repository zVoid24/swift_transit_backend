package domain

type Route struct {
	ID          int64
	Name        string
	PathGeoJSON string
	Stops       []Stop
}
