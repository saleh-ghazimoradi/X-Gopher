package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"strings"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
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
