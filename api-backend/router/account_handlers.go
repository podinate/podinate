package router

import (
	"context"
	"net/http"

	"github.com/johncave/podinate/api-backend/account"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/iam"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/responder"
	"github.com/johncave/podinate/api-backend/user"
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
		lh.Warn(ctx, "user got non-existent account", "error", apiErr, "account", requestedAccount)
		return responder.Response(http.StatusNotFound, apiErr.Error()), nil
	}

	// Check if the user can view the account
	if iam.RequestorCan(ctx, theAccount, &theAccount, account.ActionView) {
		return responder.Response(http.StatusOK, theAccount.ToAPIAccount()), nil
	}
	return responder.Response(http.StatusNotFound, "Account not found"), nil
}

// AccountPatch - Update the current account
func (s *AccountAPIService) AccountPatch(ctx context.Context, requestedAccount string, accountNew api.Account) (api.ImplResponse, error) {
	// TODO - update AccountPatch with the required logic for this service method.
	workAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, apiErr.Error()), nil
	}

	if !iam.RequestorCan(ctx, workAccount, workAccount, account.ActionUpdate) {
		return responder.Response(http.StatusNotFound, "Account not found"), nil
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
	u := iam.GetFromContext(ctx).(*user.User)
	reqID := lh.GetRequestID(ctx)

	newAcc, err := account.Create(requestedAccount, u)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(err.Code, err.Error()), nil
	}

	// Add initial policies to the account
	superAdminPolicyDocument := `
version: 2023.1
statements:
  - effect: allow
    actions: ["**"]
    resources: ["**"]`
	superAdminPolicy, err := iam.CreatePolicyForAccount(newAcc, "super-administrator", superAdminPolicyDocument, "Default policy created during initial account creation")
	err = superAdminPolicy.AttachToRequestor(u, u)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(err.Code, err.Error()), nil
	}

	lh.Log.Info("Account created", "request_id", reqID, "account", newAcc)
	return api.Response(http.StatusCreated, newAcc.ToAPIAccount()), nil
}

func (s *AccountAPIService) AccountDelete(ctx context.Context, requestedAccount string) (api.ImplResponse, error) {
	// TODO - update AccountDelete with the required logic for this service method.
	workAccount, apiErr := account.GetByID(requestedAccount)
	if apiErr.Code != http.StatusOK {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, apiErr.Error()), nil
	}

	if !iam.RequestorCan(ctx, workAccount, workAccount, account.ActionDelete) {
		return responder.Response(http.StatusNotFound, "Account not found"), nil
	}

	err := workAccount.Delete()
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusBadRequest, err.Error()), nil
	}
	return api.Response(http.StatusAccepted, workAccount.ToAPIAccount()), nil

}
