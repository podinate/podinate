package main

import (
	"context"
	"encoding/json"
	"net/http"

	account "github.com/johncave/podinate/api-backend/account"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
)

func (s *APIService) AccountGet(ctx context.Context, requestedAccount string) (api.ImplResponse, error) {
	// Get the account information from the database
	theAccount := account.Account{}
	err := theAccount.GetBySlug(requestedAccount)
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusNotFound, err.Error()), nil
	}

	return responder.Response(http.StatusOK, theAccount.ToAPIAccount()), nil
}

func (s *APIService) AccountPatch(ctx context.Context, requestedAccount string, accountNew api.Account) (api.ImplResponse, error) {
	// TODO - update AccountPatch with the required logic for this service method.

	out, _ := json.Marshal(api.Account{Id: "test", Name: "test"})
	return api.Response(http.StatusNotImplemented, string(out)), nil
}

func (s *APIService) AccountPost(ctx context.Context, requestedAccount api.Account) (api.ImplResponse, error) {
	newAcc := account.Account{Slug: requestedAccount.Id, Name: requestedAccount.Name}
	err := newAcc.ValidateNew()
	if err != nil {
		// We can pass this error directly to the API response
		return responder.Response(http.StatusBadRequest, err.Error()), nil
	}

	_, err = config.DB.Exec("INSERT INTO account(id, slug, name) VALUES(gen_random_uuid(), $1, $2)", requestedAccount.Id, requestedAccount.Name)
	// Check if insert was successful
	if err != nil {

		return responder.Response(http.StatusBadRequest, "Account ID not available"), nil
	}

	return api.Response(http.StatusCreated, newAcc.ToAPIAccount()), nil
}
