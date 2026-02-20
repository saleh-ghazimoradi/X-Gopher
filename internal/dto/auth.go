package dto

import "github.com/saleh-ghazimoradi/X-Gopher/internal/helper"

type RegisterReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResp struct {
	User         UserResp `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

func validateFirstName(v *helper.Validator, firstName string) {
	v.Check(firstName != "", "first_name", "required")
	v.Check(len(firstName) >= 2, "first_name", "length must be between 2 and 32")
	v.Check(len(firstName) <= 32, "first_name", "length must be between 2 and 32")
}

func validateLastName(v *helper.Validator, lastName string) {
	v.Check(lastName != "", "last_name", "required")
	v.Check(len(lastName) >= 2, "last_name", "length must be between 2 and 32")
	v.Check(len(lastName) <= 32, "last_name", "length must be between 2 and 32")
}

func validateEmail(v *helper.Validator, email string) {
	v.Check(email != "", "email", "required")
	v.Check(helper.Matches(email, helper.EmailRX), "email", "must be a valid email address")
}

func validatePassword(v *helper.Validator, password string) {
	v.Check(password != "", "password", "required")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters")
	v.Check(len(password) <= 72, "password", "must be at least 32 characters")
}

func ValidateRegisterReq(v *helper.Validator, req *RegisterReq) {
	validateFirstName(v, req.FirstName)
	validateLastName(v, req.LastName)
	validateEmail(v, req.Email)
	validatePassword(v, req.Password)
}

func ValidateLoginReq(v *helper.Validator, req *LoginReq) {
	validateEmail(v, req.Email)
	validatePassword(v, req.Password)
}
