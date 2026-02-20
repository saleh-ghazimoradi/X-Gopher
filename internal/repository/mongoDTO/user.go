package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	FirstName string        `bson:"first_name"`
	LastName  string        `bson:"last_name"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
	ImageUrl  string        `bson:"image_url"`
	Bio       string        `bson:"bio"`
	Followers []string      `bson:"followers"`
	Following []string      `bson:"following"`
}

func FromUserCoreToDTO(input *domain.User) (*User, error) {
	var objectId bson.ObjectID
	var err error

	if input.Id != "" {
		objectId, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
	} else {
		objectId = bson.NewObjectID()
	}

	return &User{
		Id:        objectId,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		ImageUrl:  input.ImageUrl,
		Bio:       input.Bio,
		Followers: input.Followers,
		Following: input.Following,
	}, nil
}

func FromUserDTOToCore(input *User) *domain.User {
	return &domain.User{
		Id:        input.Id.Hex(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		ImageUrl:  input.ImageUrl,
		Bio:       input.Bio,
		Followers: input.Followers,
		Following: input.Following,
	}
}
