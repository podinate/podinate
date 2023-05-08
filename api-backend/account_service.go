package main

import (
	"context"
	"net/http"

	account "github.com/johncave/podinate/api-backend/account"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
)

func (s *APIService) AccountGet(ctx context.Context, requestedAccount string) (api.ImplResponse, error) {
	// Get the account information from the database
	theAccount := account.Account{}
	err := theAccount.GetByID(requestedAccount)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, err.Error()), nil
	}

	return responder.Response(http.StatusOK, theAccount.ToAPIAccount()), nil
}

func (s *APIService) AccountPatch(ctx context.Context, requestedAccount string, accountNew api.Account) (api.ImplResponse, error) {
	// TODO - update AccountPatch with the required logic for this service method.
	workAccount := account.Account{}
	err := workAccount.GetByID(requestedAccount)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, err.Error()), nil
	}
	err = workAccount.Patch(accountNew)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusBadRequest, err.Error()), nil
	}
	return api.Response(http.StatusOK, workAccount.ToAPIAccount()), nil

}

func (s *APIService) AccountPost(ctx context.Context, requestedAccount api.Account) (api.ImplResponse, error) {
	newAcc := account.Account{ID: requestedAccount.Id, Name: requestedAccount.Name}
	err := newAcc.Register()
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(err.Code, err.Error()), nil
	}

	return api.Response(http.StatusCreated, newAcc.ToAPIAccount()), nil
}

func (s *APIService) AccountDelete(ctx context.Context, requestedAccount string) (api.ImplResponse, error) {
	// TODO - update AccountDelete with the required logic for this service method.
	workAccount := account.Account{}
	err := workAccount.GetByID(requestedAccount)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, err.Error()), nil
	}
	err = workAccount.Delete()
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusBadRequest, err.Error()), nil
	}
	return api.Response(http.StatusAccepted, workAccount.ToAPIAccount()), nil

}
