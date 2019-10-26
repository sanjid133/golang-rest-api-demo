package model

import "time"

const (
	AlgoSha256 = "sha256"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	PassHash  []byte
	PassSalt  []byte
	PassIter  int
	PassAlgo  string
	Tags      []Tag

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}
