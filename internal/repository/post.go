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

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) error
	GetPostById(ctx context.Context, id string) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) error
	ToggleLike(ctx context.Context, postId, userId string) error
	AddComment(ctx context.Context, postId, commentId string) error
	DeletePost(ctx context.Context, id string) error
	GetFeedPosts(ctx context.Context, creatorId []string) error
	SearchPosts(ctx context.Context, query string) ([]*domain.Post, error)
}

type postRepository struct {
	collection *mongo.Collection
}

func (p *postRepository) CreatePost(ctx context.Context, post *domain.Post) error {
	postDTO, err := mongoDTO.FromPostCoreToDTO(post)
	if err != nil {
		return err
	}

	res, err := p.collection.InsertOne(ctx, postDTO)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		post.Id = oid.Hex()
	}

	return nil
}

func (p *postRepository) GetPostById(ctx context.Context, id string) (*domain.Post, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid post id: %w", err)
	}

	var postDTO mongoDTO.Post
	if err := p.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&postDTO); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return mongoDTO.FromPostDTOToCore(&postDTO), nil
}

func (p *postRepository) UpdatePost(ctx context.Context, post *domain.Post) error {
	oid, err := bson.ObjectIDFromHex(post.Id)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	postDTO, err := mongoDTO.FromPostCoreToDTO(post)
	if err != nil {
		return err
	}

	_, err = p.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$set": bson.M{
			"title":         postDTO.Title,
			"message":       postDTO.Message,
			"selected_file": postDTO.SelectedFile,
		},
	})

	return err
}

func (p *postRepository) ToggleLike(ctx context.Context, postId, userId string) error {
	return nil
}

func (p *postRepository) AddComment(ctx context.Context, postId, commentId string) error {
	return nil
}

func (p *postRepository) DeletePost(ctx context.Context, id string) error {
	return nil
}

func (p *postRepository) GetFeedPosts(ctx context.Context, creatorId []string) error {
	return nil
}

func (p *postRepository) SearchPosts(ctx context.Context, query string) ([]*domain.Post, error) {
	return nil, nil
}

func NewPostRepository(database *mongo.Database, collectionName string) PostRepository {
	return &postRepository{
		collection: database.Collection(collectionName),
	}
}
