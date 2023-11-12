package router

import (
	"context"

	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
	"github.com/markbates/goth"
)

// Import the default service struct
type UserAPIService struct {
	api.UserApiService
}

// NewUserAPIService creates a new service for handling requests in the default API
func NewUserAPIService() api.UserApiServicer {
	return &UserAPIService{}
}

// Build the service methods

// UserGet - Get information about the current user
func (s *UserAPIService) UserGet(ctx context.Context, account string) (api.ImplResponse, error) {
	// TODO - update UserGet with the required logic for this service method.
	return api.Response(501, nil), nil
}

// UserLoginCompleteGet - Complete the login process
func (s *UserAPIService) UserLoginCompleteGet(ctx context.Context, token string) (api.ImplResponse, error) {
	// TODO - update UserLoginCompleteGet with the required logic for this service method.
	return api.Response(501, nil), nil
}

// UserLoginInitiateGet - Initiate the login process
func (s *UserAPIService) UserLoginInitiateGet(ctx context.Context, providerName string) (api.ImplResponse, error) {
	// TODO - update UserLoginInitiateGet with the required logic for this service method.
	//gothic.BeginAuthHandler()
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return responder.Response(400, err.Error()), nil
	}
	return api.Response(501, nil), nil
}
