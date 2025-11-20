package route

import "swift_transit/domain"

type service struct {
	repo RouteRepo
}

func NewService(repo RouteRepo) Service {
	return &service{
		repo: repo,
	}
}

func (svc *service) Create(route domain.Route) (*domain.Route, error) {

}
func (svc *service) FindAll() ([]domain.Route, error) {

}
func (svc *service) FindByID(id int64) (*domain.Route, error) {

}
