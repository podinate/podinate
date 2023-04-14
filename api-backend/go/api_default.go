/*
 * Podinate API
 *
 * The API for the simple containerisation solution Podinate
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProjectGet - Returns a list of projects.
func ProjectGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// ProjectIdGet - Get an existing project given by ID
func ProjectIdGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// ProjectIdPatch - Update an existing project
func ProjectIdPatch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// ProjectPost - Create a new project
func ProjectPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}