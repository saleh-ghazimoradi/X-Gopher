package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"strings"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	Follow(ctx context.Context, followerId, followeeId string) error
	Unfollow(ctx context.Context, followerId, followeeId string) error
}

type userRepository struct {
	collection *mongo.Collection
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	userDTO, err := mongoDTO.FromUserCoreToDTO(user)
	if err != nil {
		return err
	}

	res, err := u.collection.InsertOne(ctx, userDTO)
	if err != nil {
		if u.isDuplicateEmailError(err) {
			return ErrDuplicateEmail
		}
		return err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		user.Id = oid.Hex()
	}

	return nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userDTO mongoDTO.User

	if err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&userDTO); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return mongoDTO.FromUserDTOToCore(&userDTO), nil
}

func (u *userRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	var userDTO mongoDTO.User
	if err := u.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&userDTO); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return mongoDTO.FromUserDTOToCore(&userDTO), nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	oid, err := bson.ObjectIDFromHex(user.Id)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"image_url":  user.ImageUrl,
			"bio":        user.Bio,
		},
	}

	result, err := u.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (u *userRepository) Follow(ctx context.Context, followerId, followeeId string) error {
	followerOId, err := bson.ObjectIDFromHex(followerId)
	if err != nil {
		return fmt.Errorf("invalid follower id: %w", err)
	}

	followeeOId, err := bson.ObjectIDFromHex(followeeId)
	if err != nil {
		return fmt.Errorf("invalid followee id: %w", err)
	}

	if _, err := u.collection.UpdateOne(ctx, bson.M{"_id": followeeOId}, bson.M{
		"$addToSet": bson.M{"followers": followerId},
	}); err != nil {
		return err
	}

	if _, err := u.collection.UpdateOne(ctx, bson.M{"_id": followerOId}, bson.M{
		"$addToSet": bson.M{"following": followeeId},
	}); err != nil {
		return err
	}

	return nil
}

func (u *userRepository) Unfollow(ctx context.Context, followerId, followeeId string) error {
	followerOId, err := bson.ObjectIDFromHex(followerId)
	if err != nil {
		return fmt.Errorf("invalid follower id: %w", err)
	}

	followeeOId, err := bson.ObjectIDFromHex(followeeId)
	if err != nil {
		return fmt.Errorf("invalid followee id: %w", err)
	}

	if _, err := u.collection.UpdateOne(ctx, bson.M{"_id": followeeOId}, bson.M{
		"$pull": bson.M{"followers": followerId},
	}); err != nil {
		return err
	}

	if _, err := u.collection.UpdateOne(ctx, bson.M{"_id": followerOId}, bson.M{
		"$pull": bson.M{"following": followeeId},
	}); err != nil {
		return err
	}

	return nil
}

func (u *userRepository) isDuplicateEmailError(err error) bool {
	var we mongo.WriteException
	if errors.As(err, &we) {
		for _, e := range we.WriteErrors {
			if e.Code == 11000 && strings.Contains(e.Message, "email") {
				return true
			}
		}
	}
	return false
}

func NewUserRepository(database *mongo.Database, collectionName string) UserRepository {
	return &userRepository{
		collection: database.Collection(collectionName),
	}
}
