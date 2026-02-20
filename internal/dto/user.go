package dto

import "github.com/saleh-ghazimoradi/X-Gopher/internal/helper"

type UpdateUserReq struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	ImageUrl  *string `json:"image_url"`
	Bio       *string `json:"bio"`
}

type UserResp struct {
	Id        string   `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	ImageUrl  string   `json:"image_url"`
	Bio       string   `json:"bio"`
	Followers []string `json:"followers"`
	Following []string `json:"following"`
}

func validateImageUrl(v *helper.Validator, image string) {

}

func validateBio(v *helper.Validator, bio string) {
	v.Check(len(bio) <= 72, "bio", "bio length must not be greater than 72")
}

func ValidateUpdateUserReq(v *helper.Validator, req *UpdateUserReq) {
	if req.FirstName != nil {
		validateFirstName(v, *req.FirstName)
	}
	if req.LastName != nil {
		validateLastName(v, *req.LastName)
	}

	if req.ImageUrl != nil {
		validateImageUrl(v, *req.ImageUrl)
	}
	if req.Bio != nil {
		validateBio(v, *req.Bio)
	}
}
