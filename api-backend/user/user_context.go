package user

import "context"

// Context key for users in contexts
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func GetFromContext(ctx context.Context) *User {
	return ctx.Value(ContextKey("user")).(*User)
}
