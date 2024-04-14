# \SharedVolumeApi

All URIs are relative to *https://api.podinate.com/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProjectProjectIdSharedVolumesGet**](SharedVolumeApi.md#ProjectProjectIdSharedVolumesGet) | **Get** /project/{project_id}/shared_volumes | Get a list of shared volumes for a given project
[**ProjectProjectIdSharedVolumesPost**](SharedVolumeApi.md#ProjectProjectIdSharedVolumesPost) | **Post** /project/{project_id}/shared_volumes | Create a new shared volume
[**ProjectProjectIdSharedVolumesVolumeIdDelete**](SharedVolumeApi.md#ProjectProjectIdSharedVolumesVolumeIdDelete) | **Delete** /project/{project_id}/shared_volumes/{volume_id} | Delete a shared volume
[**ProjectProjectIdSharedVolumesVolumeIdGet**](SharedVolumeApi.md#ProjectProjectIdSharedVolumesVolumeIdGet) | **Get** /project/{project_id}/shared_volumes/{volume_id} | Get a shared volume by ID
[**ProjectProjectIdSharedVolumesVolumeIdPut**](SharedVolumeApi.md#ProjectProjectIdSharedVolumesVolumeIdPut) | **Put** /project/{project_id}/shared_volumes/{volume_id} | Update a shared volume&#39;s spec



## ProjectProjectIdSharedVolumesGet

> ProjectProjectIdSharedVolumesGet200Response ProjectProjectIdSharedVolumesGet(ctx, projectId).Account(account).Page(page).Limit(limit).Execute()

Get a list of shared volumes for a given project



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/Podinate/podinate/lib/api_client"
)

func main() {
    projectId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    page := int32(56) // int32 | The page number to return (starts at 0) (optional) (default to 0)
    limit := int32(56) // int32 | The amount of items to return per page (optional) (default to 20)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.SharedVolumeApi.ProjectProjectIdSharedVolumesGet(context.Background(), projectId).Account(account).Page(page).Limit(limit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SharedVolumeApi.ProjectProjectIdSharedVolumesGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdSharedVolumesGet`: ProjectProjectIdSharedVolumesGet200Response
    fmt.Fprintf(os.Stdout, "Response from `SharedVolumeApi.ProjectProjectIdSharedVolumesGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdSharedVolumesGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 
 **page** | **int32** | The page number to return (starts at 0) | [default to 0]
 **limit** | **int32** | The amount of items to return per page | [default to 20]

### Return type

[**ProjectProjectIdSharedVolumesGet200Response**](ProjectProjectIdSharedVolumesGet200Response.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdSharedVolumesPost

> SharedVolume ProjectProjectIdSharedVolumesPost(ctx, projectId).Account(account).SharedVolume(sharedVolume).Execute()

Create a new shared volume



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/Podinate/podinate/lib/api_client"
)

func main() {
    projectId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    sharedVolume := *openapiclient.NewSharedVolume("blog-files", int32(10)) // SharedVolume | A JSON object containing the information needed to create a new shared volume

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.SharedVolumeApi.ProjectProjectIdSharedVolumesPost(context.Background(), projectId).Account(account).SharedVolume(sharedVolume).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SharedVolumeApi.ProjectProjectIdSharedVolumesPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdSharedVolumesPost`: SharedVolume
    fmt.Fprintf(os.Stdout, "Response from `SharedVolumeApi.ProjectProjectIdSharedVolumesPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdSharedVolumesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 
 **sharedVolume** | [**SharedVolume**](SharedVolume.md) | A JSON object containing the information needed to create a new shared volume | 

### Return type

[**SharedVolume**](SharedVolume.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdSharedVolumesVolumeIdDelete

> ProjectProjectIdSharedVolumesVolumeIdDelete(ctx, projectId, volumeId).Account(account).Execute()

Delete a shared volume



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/Podinate/podinate/lib/api_client"
)

func main() {
    projectId := "hello-world" // string | 
    volumeId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdDelete(context.Background(), projectId, volumeId).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdDelete``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdSharedVolumesVolumeIdDeleteRequest struct via the builder pattern


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


## ProjectProjectIdSharedVolumesVolumeIdGet

> SharedVolume ProjectProjectIdSharedVolumesVolumeIdGet(ctx, projectId, volumeId).Account(account).Execute()

Get a shared volume by ID



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/Podinate/podinate/lib/api_client"
)

func main() {
    projectId := "hello-world" // string | 
    volumeId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdGet(context.Background(), projectId, volumeId).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdSharedVolumesVolumeIdGet`: SharedVolume
    fmt.Fprintf(os.Stdout, "Response from `SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdSharedVolumesVolumeIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 

### Return type

[**SharedVolume**](SharedVolume.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdSharedVolumesVolumeIdPut

> SharedVolume ProjectProjectIdSharedVolumesVolumeIdPut(ctx, projectId, volumeId).Account(account).SharedVolume(sharedVolume).Execute()

Update a shared volume's spec



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/Podinate/podinate/lib/api_client"
)

func main() {
    projectId := "hello-world" // string | 
    volumeId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    sharedVolume := *openapiclient.NewSharedVolume("blog-files", int32(10)) // SharedVolume | A JSON object containing the information needed to update a shared volume

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdPut(context.Background(), projectId, volumeId).Account(account).SharedVolume(sharedVolume).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdPut``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdSharedVolumesVolumeIdPut`: SharedVolume
    fmt.Fprintf(os.Stdout, "Response from `SharedVolumeApi.ProjectProjectIdSharedVolumesVolumeIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**volumeId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdSharedVolumesVolumeIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 
 **sharedVolume** | [**SharedVolume**](SharedVolume.md) | A JSON object containing the information needed to update a shared volume | 

### Return type

[**SharedVolume**](SharedVolume.md)

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

