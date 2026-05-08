package domain

import "time"

type Session struct {
	Id                 string
	HashedRefreshToken string
	UserId             int
	Device             string
	Ip                 string
	ExpiresAt          time.Time
	CreatedAt          time.Time
	IsActive           bool
}
