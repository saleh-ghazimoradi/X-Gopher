package dto

import (
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"time"
)

type CreatePostReq struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	SelectedFile string `json:"selected_file"`
}

type UpdatePostReq struct {
	Title        *string `json:"title"`
	Message      *string `json:"message"`
	SelectedFile *string `json:"selected_file"`
}

type PostResp struct {
	Id           string    `json:"id"`
	Creator      string    `json:"creator"`
	Title        string    `json:"title"`
	Message      string    `json:"message"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	SelectedFile string    `json:"selected_file"`
	Likes        []string  `json:"likes"`
	Comments     []string  `json:"comments"`
	CreatedAt    time.Time `json:"created_at"`
}

func validatePostTitle(v *helper.Validator, title string) {
	v.Check(title != "", "title", "must be provided")
	v.Check(len(title) <= 500, "title", "must not be more than 500 characters")
}

func validatePostMessage(v *helper.Validator, message string) {
	v.Check(message != "", "message", "must be provided")
	v.Check(len(message) >= 5, "message", "must be more than 5 characters")
}

func validateUpdatePostTitle(v *helper.Validator, title *string) {
	if title != nil {
		v.Check(len(*title) <= 500, "title", "must not be more than 500 characters")
	}
}

func validateUpdatePostMessage(v *helper.Validator, message *string) {
	if message != nil {
		v.Check(len(*message) >= 5, "message", "must be more than 5 characters")
	}
}

func ValidateCreatePostReq(v *helper.Validator, req *CreatePostReq) {
	validatePostTitle(v, req.Title)
	validatePostMessage(v, req.Message)
}

func ValidateUpdatePostReq(v *helper.Validator, req *UpdatePostReq) {
	validateUpdatePostTitle(v, req.Title)
	validateUpdatePostMessage(v, req.Message)
}
