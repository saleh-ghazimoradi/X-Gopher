package domain

import "time"

type RefreshToken struct {
	Id        string
	UserId    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
