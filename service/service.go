package service

import (
	"github.com/sanjid133/rest-user-store/repo"
)

type Service struct {
	usrRepo repo.User
}

func NewService(usr *repo.MgoUser) *Service  {
	return &Service{
		usrRepo: usr,
	}
}