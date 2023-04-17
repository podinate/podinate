package main

import api "github.com/johncave/podinate/api-backend/go"

type APIService struct {
	api.DefaultApiService
}

// NewAPIService creates a new service for handling requests in the default API
func NewAPIService() api.DefaultApiServicer {
	return &APIService{}
}
