/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate. Login should be performed over oauth from [auth.podinate.com](https://auth.podinate.com)
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"
	"errors"
)

// DefaultApiService is a service that implements the logic for the DefaultApiServicer
// This service should implement the business logic for every endpoint for the DefaultApi API.
// Include any external packages or services that will be required by this service.
type DefaultApiService struct {
}

// NewDefaultApiService creates a default api service
func NewDefaultApiService() DefaultApiServicer {
	return &DefaultApiService{}
}

// AccountGet - Get information about the current account.
func (s *DefaultApiService) AccountGet(ctx context.Context, account string) (ImplResponse, error) {
	// TODO - update AccountGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Account{}) or use other options such as http.Ok ...
	//return Response(200, Account{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("AccountGet method not implemented")
}

// AccountPatch - Update an existing account
func (s *DefaultApiService) AccountPatch(ctx context.Context, account string, account2 Account) (ImplResponse, error) {
	// TODO - update AccountPatch with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Account{}) or use other options such as http.Ok ...
	//return Response(200, Account{}), nil

	//TODO: Uncomment the next line to return response Response(400, Model400Error{}) or use other options such as http.Ok ...
	//return Response(400, Model400Error{}), nil

	//TODO: Uncomment the next line to return response Response(500, Model500Error{}) or use other options such as http.Ok ...
	//return Response(500, Model500Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("AccountPatch method not implemented")
}

// AccountPost - Create a new account
func (s *DefaultApiService) AccountPost(ctx context.Context, account Account) (ImplResponse, error) {
	// TODO - update AccountPost with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, Account{}) or use other options such as http.Ok ...
	//return Response(201, Account{}), nil

	//TODO: Uncomment the next line to return response Response(400, Model400Error{}) or use other options such as http.Ok ...
	//return Response(400, Model400Error{}), nil

	//TODO: Uncomment the next line to return response Response(500, Model500Error{}) or use other options such as http.Ok ...
	//return Response(500, Model500Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("AccountPost method not implemented")
}

// ProjectGet - Returns a list of projects.
func (s *DefaultApiService) ProjectGet(ctx context.Context, account string) (ImplResponse, error) {
	// TODO - update ProjectGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Project{}) or use other options such as http.Ok ...
	//return Response(200, Project{}), nil

	//TODO: Uncomment the next line to return response Response(500, Model500Error{}) or use other options such as http.Ok ...
	//return Response(500, Model500Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ProjectGet method not implemented")
}

// ProjectIdGet - Get an existing project given by ID
func (s *DefaultApiService) ProjectIdGet(ctx context.Context, id string, account string) (ImplResponse, error) {
	// TODO - update ProjectIdGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Project{}) or use other options such as http.Ok ...
	//return Response(200, Project{}), nil

	//TODO: Uncomment the next line to return response Response(400, Model400Error{}) or use other options such as http.Ok ...
	//return Response(400, Model400Error{}), nil

	//TODO: Uncomment the next line to return response Response(500, Model500Error{}) or use other options such as http.Ok ...
	//return Response(500, Model500Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ProjectIdGet method not implemented")
}

// ProjectIdPatch - Update an existing project
func (s *DefaultApiService) ProjectIdPatch(ctx context.Context, id string, account string) (ImplResponse, error) {
	// TODO - update ProjectIdPatch with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, Project{}) or use other options such as http.Ok ...
	//return Response(200, Project{}), nil

	//TODO: Uncomment the next line to return response Response(400, Model400Error{}) or use other options such as http.Ok ...
	//return Response(400, Model400Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ProjectIdPatch method not implemented")
}

// ProjectPost - Create a new project
func (s *DefaultApiService) ProjectPost(ctx context.Context, account string, project Project) (ImplResponse, error) {
	// TODO - update ProjectPost with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, Project{}) or use other options such as http.Ok ...
	//return Response(201, Project{}), nil

	//TODO: Uncomment the next line to return response Response(400, Model400Error{}) or use other options such as http.Ok ...
	//return Response(400, Model400Error{}), nil

	//TODO: Uncomment the next line to return response Response(500, Model500Error{}) or use other options such as http.Ok ...
	//return Response(500, Model500Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ProjectPost method not implemented")
}
