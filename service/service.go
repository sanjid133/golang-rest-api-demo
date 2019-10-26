package service

import (
	"github.com/sanjid133/rest-user-store/repo"
)

type Service struct {
	usrRepo repo.User
	tagRepo repo.Tag
}

func NewService(usr *repo.MgoUser, tag *repo.MgoTag) *Service {
	return &Service{
		usrRepo: usr,
		tagRepo: tag,
	}
}
