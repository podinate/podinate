package iam

import (
	"context"
	"net/http"

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

// AddRequestorToRequest - Add the user to the context
// If you need to add a new authentication mechanism, this is the place
func AddRequestorToRequest(r *http.Request) (*http.Request, error) {
	// Get the user from the database
	theUser, err := user.GetFromAPIKey(r.Header.Get("Authorization"))
	if err != nil {
		lh.Log.Errorw("Could not get requestor from API key", "error", err, "request_id", lh.GetRequestID(r.Context()))
		return r, err
	}

	// Set the user in the context
	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextKey("requestor"), theUser)

	// Add the new context to the request
	r = r.Clone(ctx)

	return r, nil
}
