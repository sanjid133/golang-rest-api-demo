package model

import "time"

type Tag struct {
	ID       string
	Tag      string
	UserID   string
	ExpireAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostTags struct {
	Tags   []string `json:"tags"`
	Expiry int32    `json:"expiry"`
}
