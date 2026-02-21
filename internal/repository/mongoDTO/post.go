package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Post struct {
	Id           bson.ObjectID `bson:"_id,omitempty"`
	Creator      string        `bson:"creator"`
	Title        string        `bson:"title"`
	Message      string        `bson:"message"`
	Name         string        `bson:"name"`
	SelectedFile string        `bson:"selected_file"`
	Likes        []string      `bson:"likes"`
	Comments     []string      `bson:"comments"`
	CreatedAt    time.Time     `bson:"created_at"`
}

func FromPostCoreToDTO(input *domain.Post) (*Post, error) {
	var objectId bson.ObjectID
	var err error

	if input.Id != "" {
		objectId, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid post id: %w", err)
		}
	} else {
		objectId = bson.NewObjectID()
	}

	return &Post{
		Id:           objectId,
		Creator:      input.Creator,
		Title:        input.Title,
		Message:      input.Message,
		Name:         input.Name,
		SelectedFile: input.SelectedFile,
		Likes:        input.Likes,
		Comments:     input.Comments,
		CreatedAt:    input.CreatedAt,
	}, nil
}

func FromPostDTOToCore(input *Post) *domain.Post {
	return &domain.Post{
		Id:           input.Id.Hex(),
		Creator:      input.Creator,
		Title:        input.Title,
		Message:      input.Message,
		Name:         input.Name,
		SelectedFile: input.SelectedFile,
		Likes:        input.Likes,
		Comments:     input.Comments,
		CreatedAt:    input.CreatedAt,
	}
}
