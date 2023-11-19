package iam

import (
	"context"

	eh "github.com/johncave/podinate/api-backend/errorhandler"
	"github.com/johncave/podinate/api-backend/user"
)

// Context key for users in contexts
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func GetFromContext(ctx context.Context) Resource {
	r, err := ctx.Value(ContextKey("requestor")).(*user.User)
	if !err {
		eh.Log.Errorw("Error getting user from context", "error", err)
		return nil
	}
	return r
}

func GetFromAuthorizationHeader(authHeader string) (Resource, error) {
	// Get the user from the database
	theUser, err := user.GetFromAPIKey(authHeader)
	if err != nil {
		eh.Log.Errorw("Error getting requestor from API key", "error", err)
		return nil, err
	}

	return theUser, nil
}
