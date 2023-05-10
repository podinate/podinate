package router

import (
	"context"
	"net/http"

	"github.com/johncave/podinate/api-backend/account"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
)

// / Import the default service
type AccountAPIService struct {
	api.AccountApiService
}

// NewAccountAPIService creates a new service for handling requests in the default API
func NewAccountAPIService() api.AccountApiServicer {
	return &AccountAPIService{}
}

/// Build the service methods

// AccountGet - Get information about the current account
func (s *AccountAPIService) AccountGet(ctx context.Context, requestedAccount string) (api.ImplResponse, error) {
	// Get the account information from the database
	theAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, apiErr.Error()), nil
	}

	return responder.Response(http.StatusOK, theAccount.ToAPIAccount()), nil
}

// AccountPatch - Update the current account
func (s *AccountAPIService) AccountPatch(ctx context.Context, requestedAccount string, accountNew api.Account) (api.ImplResponse, error) {
	// TODO - update AccountPatch with the required logic for this service method.
	workAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, apiErr.Error()), nil
	}
	apiErr = workAccount.Patch(accountNew)
	if apiErr.Code != http.StatusOK {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusBadRequest, apiErr.Error()), nil
	}
	return api.Response(http.StatusOK, workAccount.ToAPIAccount()), nil

}

// AccountPost - Request a new account
func (s *AccountAPIService) AccountPost(ctx context.Context, requestedAccount api.Account) (api.ImplResponse, error) {
	newAcc, err := account.Create(requestedAccount)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(err.Code, err.Error()), nil
	}

	return api.Response(http.StatusCreated, newAcc.ToAPIAccount()), nil
}

func (s *AccountAPIService) AccountDelete(ctx context.Context, requestedAccount string) (api.ImplResponse, error) {
	// TODO - update AccountDelete with the required logic for this service method.
	workAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr.Code != http.StatusOK {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, apiErr.Error()), nil
	}
	err := workAccount.Delete()
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusBadRequest, err.Error()), nil
	}
	return api.Response(http.StatusAccepted, workAccount.ToAPIAccount()), nil

}
