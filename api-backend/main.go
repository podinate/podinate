/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//sw "github.com/GIT_USER_ID/GIT_REPO_ID/go"
	//
	sw "github.com/johncave/podinate/api-backend/go"
)

func main() {
	log.Printf("Starting server on port 8080...")

	router := sw.NewRouter()

	log.Fatal(router.Run(":8080"))
}
