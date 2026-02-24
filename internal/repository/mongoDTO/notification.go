package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type NotificationUser struct {
	Name   string `bson:"name"`
	Avatar string `bson:"avatar"`
}

type Notification struct {
	Id               bson.ObjectID    `bson:"_id,omitempty"`
	SenderId         bson.ObjectID    `bson:"sender_id"`
	ReceiverId       bson.ObjectID    `bson:"receiver_id"`
	TargetId         bson.ObjectID    `bson:"target_id"`
	Details          string           `bson:"details"`
	IsRead           bool             `bson:"is_read"`
	CreatedAt        time.Time        `bson:"created_at"`
	NotificationUser NotificationUser `bson:"notification_user"`
}

func FromNotificationDTOToCore(input *Notification) *domain.Notification {
	return &domain.Notification{
		Id:         input.Id.Hex(),
		SenderId:   input.SenderId.Hex(),
		ReceiverId: input.ReceiverId.Hex(),
		Details:    input.Details,
		IsRead:     input.IsRead,
		CreatedAt:  input.CreatedAt,
		NotificationUser: domain.NotificationUser{
			Name:   input.NotificationUser.Name,
			Avatar: input.NotificationUser.Avatar,
		},
	}
}

func FromNotificationCoreTODTO(input *domain.Notification) (*Notification, error) {
	senderOID, err := bson.ObjectIDFromHex(input.SenderId)
	if err != nil {
		return nil, fmt.Errorf("invalid sender id: %w", err)
	}
	receiverOID, err := bson.ObjectIDFromHex(input.ReceiverId)
	if err != nil {
		return nil, fmt.Errorf("invalid receiver id: %w", err)
	}

	var targetOID bson.ObjectID
	if input.TargetId != "" {
		targetOID, err = bson.ObjectIDFromHex(input.TargetId)
		if err != nil {
			return nil, fmt.Errorf("invalid target id: %w", err)
		}
	}

	var idOID bson.ObjectID
	if input.Id != "" {
		idOID, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid notification id: %w", err)
		}
	} else {
		idOID = bson.NewObjectID()
	}

	return &Notification{
		Id:         idOID,
		SenderId:   senderOID,
		ReceiverId: receiverOID,
		TargetId:   targetOID,
		Details:    input.Details,
		IsRead:     input.IsRead,
		CreatedAt:  input.CreatedAt,
		NotificationUser: NotificationUser{
			Name:   input.NotificationUser.Name,
			Avatar: input.NotificationUser.Avatar,
		},
	}, nil
}
