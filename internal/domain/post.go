package domain

import "time"

type Post struct {
	Id           string
	Creator      string
	Title        string
	Message      string
	Name         string
	SelectedFile string
	Likes        []string
	Comments     []string
	CreatedAt    time.Time
}
