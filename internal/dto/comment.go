package dto

import "github.com/saleh-ghazimoradi/X-Gopher/internal/helper"

type CommentReq struct {
	Value string `json:"value"`
}

func validateComment(v *helper.Validator, value string) {
	v.Check(value != "", "value", "must be provided")
	v.Check(len(value) >= 1, "value", "must be at least 1 character")
	v.Check(len(value) <= 500, "value", "must not exceed 500 characters")
}

func ValidateComment(v *helper.Validator, req *CommentReq) {
	validateComment(v, req.Value)
}
