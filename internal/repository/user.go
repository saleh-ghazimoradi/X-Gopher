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
	GetUsersByIds(ctx context.Context, ids []string) ([]*domain.User, error)
	GetUsersBySearch(ctx context.Context, query string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	Follow(ctx context.Context, followerId, followeeId string) error
	Unfollow(ctx context.Context, followerId, followeeId string) error
	DeleteUser(ctx context.Context, id string) error
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

func (u *userRepository) GetUsersByIds(ctx context.Context, ids []string) ([]*domain.User, error) {
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}

	objectIDs := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		oid, err := bson.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, oid)
	}

	cursor, err := u.collection.Find(ctx, bson.M{
		"_id": bson.M{"$in": objectIDs},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var usersDTO []mongoDTO.User
	if err := cursor.All(ctx, &usersDTO); err != nil {
		return nil, err
	}

	users := make([]*domain.User, 0, len(usersDTO))
	for _, dto := range usersDTO {
		users = append(users, mongoDTO.FromUserDTOToCore(&dto))
	}

	return users, nil
}

func (u *userRepository) GetUsersBySearch(ctx context.Context, query string) ([]*domain.User, error) {
	if query == "" {
		return nil, nil
	}

	filter := bson.M{
		"$or": []bson.M{
			{"first_name": bson.M{"$regex": query, "$options": "i"}},
			{"last_name": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var usersDTO []mongoDTO.User
	if err := cursor.All(ctx, &usersDTO); err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(usersDTO))
	for i, dto := range usersDTO {
		users[i] = mongoDTO.FromUserDTOToCore(&dto)
	}

	return users, nil

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

func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", err)
	}

	if _, err := u.collection.DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return ErrRecordNotFound
		default:
			return err
		}
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
