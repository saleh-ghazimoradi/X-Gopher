package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Comment struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	PostId    bson.ObjectID `bson:"post_id"`
	UserId    bson.ObjectID `bson:"user_id"`
	Value     string        `bson:"value"`
	CreatedAt time.Time     `bson:"created_at"`
}

func FromCoreCommentToDTO(input *domain.Comment) (*Comment, error) {
	postDTO, err := bson.ObjectIDFromHex(input.PostId)
	if err != nil {
		return nil, fmt.Errorf("invalid post id: %w", err)
	}

	userDTO, err := bson.ObjectIDFromHex(input.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	var idOID bson.ObjectID
	if input.Id != "" {
		idOID, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, err
		}
	} else {
		idOID = bson.NewObjectID()
	}

	return &Comment{
		Id:        idOID,
		PostId:    postDTO,
		UserId:    userDTO,
		Value:     input.Value,
		CreatedAt: input.CreatedAt,
	}, nil
}

func FromCommentDTOToCore(input *Comment) *domain.Comment {
	return &domain.Comment{
		Id:        input.Id.Hex(),
		PostId:    input.PostId.Hex(),
		UserId:    input.UserId.Hex(),
		Value:     input.Value,
		CreatedAt: input.CreatedAt,
	}
}
