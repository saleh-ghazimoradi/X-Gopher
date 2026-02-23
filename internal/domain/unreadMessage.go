package domain

type UnreadMessage struct {
	Id                  string
	SenderId            string
	ReceiverId          string
	NumOfUnreadMessages int
	IsRead              bool
}
