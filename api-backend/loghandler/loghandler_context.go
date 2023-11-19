package errorhandler

import (
	"context"

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
	r, err := ctx.Value(ContextKey("request-id")).(RequestID)
	if !err {
		Log.Errorw("Error getting request ID from context", "error", err)
		return ""
	}
	return r
}

func NewRequestID() RequestID {
	return "pr-" + RequestID(uuid.New().String())
}
