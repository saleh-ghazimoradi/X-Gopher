package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *domain.Comment) error
	DeleteComment(ctx context.Context, id string) error
}

type commentRepository struct {
	collection *mongo.Collection
}

func (c *commentRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	commentDTO, err := mongoDTO.FromCoreCommentToDTO(comment)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}

	res, err := c.collection.InsertOne(ctx, commentDTO)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		comment.Id = oid.Hex()
	}

	return nil
}

func (c *commentRepository) DeleteComment(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	_, err = c.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ErrRecordNotFound
		default:
			return fmt.Errorf("faield to delete comment: %w", err)
		}
	}

	return nil
}

func NewCommentRepository(database *mongo.Database, collectionName string) CommentRepository {
	return &commentRepository{
		collection: database.Collection(collectionName),
	}
}
