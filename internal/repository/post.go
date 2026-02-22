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

type PostRepository interface {
	CreatePost(ctx context.Context, post *domain.Post) error
	GetPostById(ctx context.Context, id string) (*domain.Post, error)
	UpdatePost(ctx context.Context, post *domain.Post) error
	ToggleLike(ctx context.Context, postId, userId string) error
	AddComment(ctx context.Context, postId, commentId string) error
	DeletePost(ctx context.Context, id string) error
	GetFeedPosts(ctx context.Context, creatorIds []string, page, limit int) ([]*domain.Post, int64, error)
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
	oid, err := bson.ObjectIDFromHex(postId)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}
	_, err = p.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$addToSet": bson.M{"likes": userId},
	})
	return nil
}

func (p *postRepository) AddComment(ctx context.Context, postId, commentId string) error {
	oid, err := bson.ObjectIDFromHex(postId)
	if err != nil {
		return fmt.Errorf("invalid post id: %w", err)
	}

	_, err = p.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
		"$push": bson.M{
			"comments": commentId,
		},
	})

	return err
}

func (p *postRepository) DeletePost(ctx context.Context, id string) error {
	return nil
}

func (p *postRepository) GetFeedPosts(ctx context.Context, creatorIds []string, page, limit int) ([]*domain.Post, int64, error) {
	if len(creatorIds) == 0 {
		return nil, 0, nil
	}

	filter := bson.M{"creator": bson.M{"$in": creatorIds}}

	total, err := p.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"_id": -1}).
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit))

	cursor, err := p.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var postsDTO []mongoDTO.Post
	if err := cursor.All(ctx, &postsDTO); err != nil {
		return nil, 0, err
	}

	posts := make([]*domain.Post, len(postsDTO))
	for i, dto := range postsDTO {
		posts[i] = mongoDTO.FromPostDTOToCore(&dto)
	}

	return posts, total, nil
}

func (p *postRepository) SearchPosts(ctx context.Context, query string) ([]*domain.Post, error) {
	if query == "" {
		return nil, nil
	}

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"message": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := p.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var postsDTO []mongoDTO.Post
	if err := cursor.All(ctx, &postsDTO); err != nil {
		return nil, err
	}

	posts := make([]*domain.Post, len(postsDTO))
	for i, dto := range postsDTO {
		posts[i] = mongoDTO.FromPostDTOToCore(&dto)
	}

	return posts, nil
}

func NewPostRepository(database *mongo.Database, collectionName string) PostRepository {
	return &postRepository{
		collection: database.Collection(collectionName),
	}
}
