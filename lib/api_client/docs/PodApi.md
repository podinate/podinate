# \PodApi

All URIs are relative to *https://api.podinate.com/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProjectProjectIdPodGet**](PodApi.md#ProjectProjectIdPodGet) | **Get** /project/{project_id}/pod | Get a list of pods for a given project
[**ProjectProjectIdPodPodIdDelete**](PodApi.md#ProjectProjectIdPodPodIdDelete) | **Delete** /project/{project_id}/pod/{pod_id} | Delete a pod
[**ProjectProjectIdPodPodIdGet**](PodApi.md#ProjectProjectIdPodPodIdGet) | **Get** /project/{project_id}/pod/{pod_id} | Get a pod by ID
[**ProjectProjectIdPodPodIdPatch**](PodApi.md#ProjectProjectIdPodPodIdPatch) | **Patch** /project/{project_id}/pod/{pod_id} | Update a pod
[**ProjectProjectIdPodPost**](PodApi.md#ProjectProjectIdPodPost) | **Post** /project/{project_id}/pod | Create a new pod



## ProjectProjectIdPodGet

> ProjectProjectIdPodGet200Response ProjectProjectIdPodGet(ctx, projectId).Account(account).Page(page).Limit(limit).Execute()

Get a list of pods for a given project



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    projectId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    page := int32(56) // int32 | The page number to return (starts at 0) (optional) (default to 0)
    limit := int32(56) // int32 | The amount of items to return per page (optional) (default to 20)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodGet(context.Background(), projectId).Account(account).Page(page).Limit(limit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodGet`: ProjectProjectIdPodGet200Response
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 
 **page** | **int32** | The page number to return (starts at 0) | [default to 0]
 **limit** | **int32** | The amount of items to return per page | [default to 20]

### Return type

[**ProjectProjectIdPodGet200Response**](ProjectProjectIdPodGet200Response.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdPodPodIdDelete

> ProjectProjectIdPodPodIdDelete(ctx, projectId, podId).Account(account).Execute()

Delete a pod



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    projectId := "hello-world" // string | 
    podId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.PodApi.ProjectProjectIdPodPodIdDelete(context.Background(), projectId, podId).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPodIdDelete``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**podId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPodIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 

### Return type

 (empty response body)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdPodPodIdGet

> Pod ProjectProjectIdPodPodIdGet(ctx, projectId, podId).Account(account).Execute()

Get a pod by ID



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    projectId := "hello-world" // string | 
    podId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodPodIdGet(context.Background(), projectId, podId).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPodIdGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodPodIdGet`: Pod
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodPodIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**podId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPodIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 

### Return type

[**Pod**](Pod.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdPodPodIdPatch

> Pod ProjectProjectIdPodPodIdPatch(ctx, projectId, podId).Account(account).Pod(pod).Execute()

Update a pod



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    projectId := "hello-world" // string | 
    podId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    pod := *openapiclient.NewPod() // Pod | A JSON object containing the information needed to update a pod

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodPodIdPatch(context.Background(), projectId, podId).Account(account).Pod(pod).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPodIdPatch``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodPodIdPatch`: Pod
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodPodIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**podId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPodIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 
 **pod** | [**Pod**](Pod.md) | A JSON object containing the information needed to update a pod | 

### Return type

[**Pod**](Pod.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdPodPost

> Pod ProjectProjectIdPodPost(ctx, projectId).Account(account).Pod(pod).Execute()

Create a new pod



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
    projectId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    pod := *openapiclient.NewPod() // Pod | A JSON object containing the information needed to create a new pod

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodPost(context.Background(), projectId).Account(account).Pod(pod).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodPost`: Pod
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 
 **pod** | [**Pod**](Pod.md) | A JSON object containing the information needed to create a new pod | 

### Return type

[**Pod**](Pod.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

