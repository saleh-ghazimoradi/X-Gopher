package repository

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
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

func NewMessageRepository(database *mongo.Database, collectionName string) MessageRepository {
	return &messageRepository{
		collection: database.Collection(collectionName),
	}
}
