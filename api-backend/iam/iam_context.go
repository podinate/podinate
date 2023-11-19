package iam

import (
	"context"

	"github.com/johncave/podinate/api-backend/user"

	lh "github.com/johncave/podinate/api-backend/loghandler"
)

// Context key for users in contexts
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func GetFromContext(ctx context.Context) Resource {
	r, err := ctx.Value(ContextKey("requestor")).(*user.User)
	if !err {
		lh.Log.Errorw("Error getting user from context", "error", err)
		return nil
	}
	return r
}

func GetFromAuthorizationHeader(authHeader string) (Resource, error) {
	// Get the user from the database
	theUser, err := user.GetFromAPIKey(authHeader)
	if err != nil {
		lh.Log.Errorw("Error getting requestor from API key", "error", err)
		return nil, err
	}

	return theUser, nil
}
