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

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	MarkAsRead(ctx context.Context, userId string) error
	GetByUserId(ctx context.Context, userId string) ([]*domain.Notification, error)
}

type notificationRepository struct {
	collection *mongo.Collection
}

func (n *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	dto, err := mongoDTO.FromNotificationCoreTODTO(notification)
	if err != nil {
		return err
	}

	res, err := n.collection.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		notification.Id = oid.Hex()
	}
	return nil
}

func (n *notificationRepository) MarkAsRead(ctx context.Context, userID string) error {
	oid, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	_, err = n.collection.UpdateMany(ctx,
		bson.M{"receiver_id": oid},
		bson.M{"$set": bson.M{"is_read": true}},
	)
	return err
}

func (n *notificationRepository) GetByUserId(ctx context.Context, userId string) ([]*domain.Notification, error) {
	oid, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := n.collection.Find(ctx, bson.M{"receiver_id": oid}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dtos []mongoDTO.Notification
	if err := cursor.All(ctx, &dtos); err != nil {
		return nil, err
	}

	notifications := make([]*domain.Notification, len(dtos))
	for i, dto := range dtos {
		notifications[i] = mongoDTO.FromNotificationDTOToCore(&dto)
	}

	return notifications, nil
}

func NewNotificationRepository(database *mongo.Database, collectionName string) NotificationRepository {
	return &notificationRepository{
		collection: database.Collection(collectionName),
	}
}
