package model

import "time"

type User struct {
	ID string
	FirstName string
	LastName string
	Password string
	Tags []Tag

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tag struct {
	Tag string
	ExpireAt time.Time
}

type PostUser struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}