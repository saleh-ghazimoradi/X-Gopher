package domain

import "time"

type Comment struct {
	Id        string
	PostId    string
	UserId    string
	Value     string
	CreatedAt time.Time
}
