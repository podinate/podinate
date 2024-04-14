package lh

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// Context key for users in contexts
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type RequestID string

func (c RequestID) String() string {
	return string(c)
}

func GetRequestID(ctx context.Context) RequestID {

	id := ctx.Value(ContextKey("request-id"))

	switch id.(type) {
	case RequestID:
		return id.(RequestID)
	default:
		log.Printf("Error getting request ID from context, type was %T\n", id)
		return ""
	}
}

func NewRequestID() RequestID {
	return "pr-" + RequestID(uuid.New().String())
}

func TestContext() context.Context {
	ctx := context.Background()
	rid := NewRequestID()
	log.Printf("Using request ID: %s\n", rid)
	return context.WithValue(ctx, ContextKey("request-id"), rid)
}
