package mongoDTO

import (
	"github.com/saleh-ghazimoradi/X-Gopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type RefreshToken struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserId    bson.ObjectID `bson:"user_id"`
	Token     string        `bson:"token"`
	ExpiresAt time.Time     `bson:"expires_at"`
	CreatedAt time.Time     `bson:"created_at"`
}

func FromCoreRefreshTokenToDTO(input *domain.RefreshToken) (*RefreshToken, error) {
	userOID, err := bson.ObjectIDFromHex(input.UserId)
	if err != nil {
		return nil, err
	}

	var tokenOID bson.ObjectID
	if input.Id != "" {
		tokenOID, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, err
		}
	}

	return &RefreshToken{
		Id:        tokenOID,
		UserId:    userOID,
		Token:     input.Token,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: input.CreatedAt,
	}, nil
}

func FromRefreshTokenDTOToCore(input *RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		Id:        input.Id.Hex(),
		UserId:    input.UserId.Hex(),
		Token:     input.Token,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: input.CreatedAt,
	}
}
