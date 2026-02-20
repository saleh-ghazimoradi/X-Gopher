package dto

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
