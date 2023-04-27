package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
	"github.com/johncave/podinate/api-backend/responder"
)

func (s *APIService) AccountGet(ctx context.Context, account string) (api.ImplResponse, error) {
	// TODO - update AccountGet with the required logic for this service method.

	out, _ := json.Marshal(api.Account{Id: "test", Name: "test"})
	return api.Response(http.StatusNotImplemented, string(out)), nil
}

func (s *APIService) AccountPatch(ctx context.Context, account string, accountNew api.Account) (api.ImplResponse, error) {
	// TODO - update AccountPatch with the required logic for this service method.

	out, _ := json.Marshal(api.Account{Id: "test", Name: "test"})
	return api.Response(http.StatusNotImplemented, string(out)), nil
}

func (s *APIService) AccountPost(ctx context.Context, account api.Account) (api.ImplResponse, error) {
	// Create an account in DB

	//defer stmt.Close()
	//row := stmt.QueryRow(account.Id, account.Name)
	_, err := config.DB.Exec("INSERT INTO account(id, slug, name) VALUES(gen_random_uuid(), $1, $2)", account.Id, account.Name)
	// Check if insert was successful
	if err != nil {
		log.Println("Peepeepoopoo error happened")
		// out := api.Model500Error{Code: http.StatusInternalServerError, Message: "Failed to create account"}
		// return responder.Response(http.StatusInternalServerError, out), err
		return responder.Response(http.StatusInternalServerError, "For fuck sake"), nil
	}

	return api.Response(http.StatusCreated, account), nil
}
