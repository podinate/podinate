package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	api "github.com/johncave/podinate/api-backend/go"
)

// ProjectGet - Returns a list of projects.
func (s *APIService) ProjectGet(ctx context.Context, account string) (api.ImplResponse, error) {
	// TODO - update ProjectGet with the required logic for this service method.

	out, _ := json.Marshal(api.App{Id: "podinate-blog", Name: "Podinate Blog", Image: "wordpress", Tag: "latest"})
	log.Println("Fuck you")
	return api.Response(http.StatusNotImplemented, string(out)), nil
}

func (s *APIService) ProjectIdGet(ctx context.Context, id string, account string) (api.ImplResponse, error) {
	// TODO - update ProjectIdGet with the required logic for this service method.

	out, _ := json.Marshal(api.App{Id: id, Name: "test", Image: "test", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, out), errors.New("ProjectIdGet method not implemented")
}

func (s *APIService) ProjectIdPatch(ctx context.Context, id string, account string) (api.ImplResponse, error) {
	// TODO - update ProjectIdPatch with the required logic for this service method.

	out, _ := json.Marshal(api.App{Id: id, Name: "test", Image: "test", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, out), errors.New("ProjectIdPatch method not implemented")
}

func (s *APIService) ProjectPost(ctx context.Context, account string, app api.App) (api.ImplResponse, error) {
	// TODO - update ProjectPost with the required logic for this service method.

	out, _ := json.Marshal(api.App{Id: "elephant", Name: "test", Image: "test", Tag: "latest"})
	return api.Response(http.StatusNotImplemented, out), errors.New("ProjectPost method not implemented")
}
