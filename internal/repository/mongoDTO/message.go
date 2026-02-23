package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Message struct {
	Id       bson.ObjectID `bson:"_id,omitempty"`
	Content  string        `bson:"content"`
	Sender   string        `bson:"sender"`
	Receiver string        `bson:"receiver"`
}

func FromMessageCoreToDTO(input *domain.Message) (*Message, error) {
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

	return &Message{
		Id:       objectId,
		Content:  input.Content,
		Sender:   input.Sender,
		Receiver: input.Receiver,
	}, nil
}

func FromMessageDTOToCore(input *Message) *domain.Message {
	return &domain.Message{
		Id:       input.Id.Hex(),
		Content:  input.Content,
		Sender:   input.Sender,
		Receiver: input.Receiver,
	}
}
