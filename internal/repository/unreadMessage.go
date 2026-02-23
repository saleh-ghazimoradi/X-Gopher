package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UnreadMessageRepository interface {
	IncrementUnread(ctx context.Context, senderId, receiverId string) error
}

type unreadMessageRepository struct {
	collection *mongo.Collection
}

func (u *unreadMessageRepository) IncrementUnread(ctx context.Context, senderId, receiverId string) error {
	{

		filter := bson.M{"sender_id": senderId, "receiver_id": receiverId}
		update := bson.M{
			"$inc": bson.M{"num_of_unread_messages": 1},
			"$set": bson.M{"is_read": false},
		}

		opts := options.FindOneAndUpdate().
			SetUpsert(true).
			SetReturnDocument(options.After)

		err := u.collection.FindOneAndUpdate(ctx, filter, update, opts).
			Decode(new(domain.UnreadMessage))

		return err
	}
}

func NewUnreadMessageRepository(database *mongo.Database, collectionName string) UnreadMessageRepository {
	return &unreadMessageRepository{
		collection: database.Collection(collectionName),
	}
}
