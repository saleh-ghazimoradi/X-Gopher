package dto

import "github.com/saleh-ghazimoradi/X-Gopher/internal/helper"

type MessageReq struct {
	Content  string `json:"content"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type MessageResp struct {
	Id string `json:"id"`
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
