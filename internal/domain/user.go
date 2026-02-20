package domain

type User struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	ImageUrl  string
	Bio       string
	Followers []string
	Following []string
}
