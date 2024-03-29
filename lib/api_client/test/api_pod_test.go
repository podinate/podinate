/*
Podinate API

Testing PodApiService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package api_client

import (
	"context"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_api_client_PodApiService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test PodApiService ProjectProjectIdPodGet", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string

		resp, httpRes, err := apiClient.PodApi.ProjectProjectIdPodGet(context.Background(), projectId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test PodApiService ProjectProjectIdPodPodIdDelete", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string
		var podId string

		httpRes, err := apiClient.PodApi.ProjectProjectIdPodPodIdDelete(context.Background(), projectId, podId).Execute()

		require.Nil(t, err)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test PodApiService ProjectProjectIdPodPodIdExecPost", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string
		var podId string

		resp, httpRes, err := apiClient.PodApi.ProjectProjectIdPodPodIdExecPost(context.Background(), projectId, podId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test PodApiService ProjectProjectIdPodPodIdGet", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string
		var podId string

		resp, httpRes, err := apiClient.PodApi.ProjectProjectIdPodPodIdGet(context.Background(), projectId, podId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test PodApiService ProjectProjectIdPodPodIdLogsGet", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string
		var podId string

		resp, httpRes, err := apiClient.PodApi.ProjectProjectIdPodPodIdLogsGet(context.Background(), projectId, podId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test PodApiService ProjectProjectIdPodPodIdPut", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string
		var podId string

		resp, httpRes, err := apiClient.PodApi.ProjectProjectIdPodPodIdPut(context.Background(), projectId, podId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test PodApiService ProjectProjectIdPodPost", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		var projectId string

		resp, httpRes, err := apiClient.PodApi.ProjectProjectIdPodPost(context.Background(), projectId).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}
