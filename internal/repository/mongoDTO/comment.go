package mongoDTO

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Comment struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	PostId    bson.ObjectID `bson:"post_id"`
	UserId    bson.ObjectID `bson:"user_id"`
	Value     string        `bson:"value"`
	CreatedAt time.Time     `bson:"created_at"`
}
