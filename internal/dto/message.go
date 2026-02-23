package dto

import "github.com/saleh-ghazimoradi/X-Gopher/internal/helper"

type MessageReq struct {
	Content  string `json:"content"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type GetMessagesQuery struct {
	Sender   string
	Receiver string
	Page     int
}

type UnreadConversation struct {
	Id                  string `json:"id"`
	SenderId            string `json:"sender_id"`
	ReceiverId          string `json:"receiver_id"`
	NumOfUnreadMessages int    `json:"num_of_unread_messages"`
	IsRead              bool   `json:"is_read"`
}

type UnreadSummaryResp struct {
	Conversations []UnreadConversation `json:"conversations"`
	Total         int                  `json:"total"`
}

type MessageResp struct {
	Id       string `json:"id"`
	Content  string `json:"content"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

func validateContent(v *helper.Validator, content string) {
	v.Check(content != "", "content", "must not be empty")
}

func validateSender(v *helper.Validator, sender string) {
	v.Check(sender != "", "sender", "must not be empty")
}

func validateReceiver(v *helper.Validator, receiver string) {
	v.Check(receiver != "", "receiver", "must not be empty")
}

func ValidateMessageReq(v *helper.Validator, req *MessageReq) {
	validateContent(v, req.Content)
	validateSender(v, req.Sender)
	validateReceiver(v, req.Receiver)
}

func ValidateGetMessagesQuery(v *helper.Validator, q *GetMessagesQuery) {
	v.Check(q.Sender != "", "sender", "must not be empty")
	v.Check(q.Receiver != "", "receiver", "must not be empty")
	v.Check(q.Page >= 0, "page", "must be >= 0")
}
