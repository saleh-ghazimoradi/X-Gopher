package utils

import "context"

type ContextKey string

const (
	UserIdKey ContextKey = "user_id"
)

func WithUserId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, UserIdKey, id)
}

func UserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(UserIdKey).(string)
	return userId, ok
}
