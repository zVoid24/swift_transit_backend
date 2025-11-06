package user

import (
	"swift_transit/domain"
)

type service struct {
	userRepo UserRepo
}

func NewService(usrRepo UserRepo) Service {
	return &service{
		userRepo: usrRepo,
	}
}

func (svc *service) Create(user domain.User) (*domain.User, error) {
	usr,err:=svc.userRepo.Create(user)
	if err != nil{
		return nil,err
	}
	if usr == nil{
		return nil,nil
	}
	return usr,nil
}
func (svc *service) Find(username string, password string) (*domain.User, error) {
	usr,err:=svc.userRepo.Find(username,password)
	if err!=nil{
		return nil,err
	}
	if usr == nil{
		return nil,nil
	}
	return usr,nil
}
