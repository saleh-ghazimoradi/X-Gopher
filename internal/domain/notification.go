package domain

import "time"

type NotificationUser struct {
	Name   string
	Avatar string
}

type Notification struct {
	Id               string
	Details          string
	SenderId         string
	ReceiverId       string
	TargetId         string
	IsRead           bool
	CreatedAt        time.Time
	NotificationUser NotificationUser
}
