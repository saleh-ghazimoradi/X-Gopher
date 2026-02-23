package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UnreadMessageRepository interface {
	IncrementUnread(ctx context.Context, senderId, receiverId string) error
	GetUnreadByReceiver(ctx context.Context, receiverId string) ([]*domain.UnreadMessage, error)
	MarkAsRead(ctx context.Context, receiverId, senderId string) error
}

type unreadMessageRepository struct {
	collection *mongo.Collection
}

func (u *unreadMessageRepository) IncrementUnread(ctx context.Context, senderId, receiverId string) error {
	senderOID, err := bson.ObjectIDFromHex(senderId)
	if err != nil {
		return fmt.Errorf("invalid sender id: %w", err)
	}

	receiverOID, err := bson.ObjectIDFromHex(receiverId)
	if err != nil {
		return fmt.Errorf("invalid receiver id: %w", err)
	}

	filter := bson.M{
		"sender_id":   senderOID,
		"receiver_id": receiverOID,
	}

	update := bson.M{
		"$inc": bson.M{"num_of_unread_messages": 1},
		"$set": bson.M{"is_read": false},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	err = u.collection.FindOneAndUpdate(ctx, filter, update, opts).
		Decode(new(domain.UnreadMessage))

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	return nil
}

func (u *unreadMessageRepository) GetUnreadByReceiver(ctx context.Context, receiverId string) ([]*domain.UnreadMessage, error) {
	receiverOID, err := bson.ObjectIDFromHex(receiverId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"receiver_id": receiverOID,
		"is_read":     false,
	}

	cursor, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.UnreadMessage

	for cursor.Next(ctx) {
		var dtoMsg mongoDTO.UnreadMessage
		if err := cursor.Decode(&dtoMsg); err != nil {
			return nil, err
		}
		results = append(results, mongoDTO.FromUnreadMessageDTOToCore(&dtoMsg))
	}

	return results, nil
}

func (u *unreadMessageRepository) MarkAsRead(ctx context.Context, receiverId, senderId string) error {
	receiverOID, err := bson.ObjectIDFromHex(receiverId)
	if err != nil {
		return err
	}

	senderOID, err := bson.ObjectIDFromHex(senderId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"receiver_id": receiverOID,
		"sender_id":   senderOID,
	}

	update := bson.M{
		"$set": bson.M{
			"is_read":                true,
			"num_of_unread_messages": 0,
		},
	}

	_, err = u.collection.UpdateOne(ctx, filter, update)
	return err
}

func NewUnreadMessageRepository(database *mongo.Database, collectionName string) UnreadMessageRepository {
	return &unreadMessageRepository{
		collection: database.Collection(collectionName),
	}
}
