package repository

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetMessagesBetween(ctx context.Context, user1, user2 string, skip, limit int64) ([]*domain.Message, error)
}

type messageRepository struct {
	collection *mongo.Collection
}

func (m *messageRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	messageDTO, err := mongoDTO.FromMessageCoreToDTO(message)
	if err != nil {
		return fmt.Errorf("failed to convert message dto: %w", err)
	}

	res, err := m.collection.InsertOne(ctx, messageDTO)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		message.Id = oid.Hex()
	}

	return nil
}

func (m *messageRepository) GetMessagesBetween(ctx context.Context, user1, user2 string, skip, limit int64,
) ([]*domain.Message, error) {

	filter := bson.M{
		"$or": []bson.M{
			{"sender": user1, "receiver": user2},
			{"sender": user2, "receiver": user1},
		},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*domain.Message

	for cursor.Next(ctx) {
		var dtoMsg mongoDTO.Message
		if err := cursor.Decode(&dtoMsg); err != nil {
			return nil, err
		}
		results = append(results, mongoDTO.FromMessageDTOToCore(&dtoMsg))
	}

	return results, nil
}

func NewMessageRepository(database *mongo.Database, collectionName string) MessageRepository {
	return &messageRepository{
		collection: database.Collection(collectionName),
	}
}
