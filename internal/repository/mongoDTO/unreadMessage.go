package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UnreadMessage struct {
	Id                  bson.ObjectID `bson:"_id,omitempty"`
	SenderId            bson.ObjectID `bson:"sender_id"`
	ReceiverId          bson.ObjectID `bson:"receiver_id"`
	NumOfUnreadMessages int           `bson:"num_of_unread_messages"`
	IsRead              bool          `bson:"is_read"`
}

func FromUnreadMessageCoreToDTO(input *domain.UnreadMessage) (*UnreadMessage, error) {
	var objectId bson.ObjectID
	var err error

	if input.SenderId != "" {
		objectId, err = bson.ObjectIDFromHex(input.SenderId)
		if err != nil {
			return nil, fmt.Errorf("invalid unread message id")
		}
	} else {
		objectId = bson.NewObjectID()
	}

	senderDTO, err := bson.ObjectIDFromHex(input.ReceiverId)
	if err != nil {
		return nil, fmt.Errorf("invalid sender Id")
	}

	receiverDTO, err := bson.ObjectIDFromHex(input.ReceiverId)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver Id")
	}

	return &UnreadMessage{
		Id:                  objectId,
		SenderId:            senderDTO,
		ReceiverId:          receiverDTO,
		NumOfUnreadMessages: input.NumOfUnreadMessages,
		IsRead:              input.IsRead,
	}, nil
}

func FromUnreadMessageDTOToCore(input *UnreadMessage) *domain.UnreadMessage {
	return &domain.UnreadMessage{
		Id:                  input.Id.Hex(),
		SenderId:            input.SenderId.Hex(),
		ReceiverId:          input.ReceiverId.Hex(),
		NumOfUnreadMessages: input.NumOfUnreadMessages,
		IsRead:              input.IsRead,
	}
}
