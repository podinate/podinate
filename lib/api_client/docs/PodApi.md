# \PodApi

All URIs are relative to *https://api.podinate.com/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProjectProjectIdPodGet**](PodApi.md#ProjectProjectIdPodGet) | **Get** /project/{project_id}/pod | Get a list of pods for a given project
[**ProjectProjectIdPodPodIdDelete**](PodApi.md#ProjectProjectIdPodPodIdDelete) | **Delete** /project/{project_id}/pod/{pod_id} | Delete a pod
[**ProjectProjectIdPodPodIdExecPost**](PodApi.md#ProjectProjectIdPodPodIdExecPost) | **Post** /project/{project_id}/pod/{pod_id}/exec | Execute a command in a pod
[**ProjectProjectIdPodPodIdGet**](PodApi.md#ProjectProjectIdPodPodIdGet) | **Get** /project/{project_id}/pod/{pod_id} | Get a pod by ID
[**ProjectProjectIdPodPodIdLogsGet**](PodApi.md#ProjectProjectIdPodPodIdLogsGet) | **Get** /project/{project_id}/pod/{pod_id}/logs | Get the logs for a pod
[**ProjectProjectIdPodPodIdPut**](PodApi.md#ProjectProjectIdPodPodIdPut) | **Put** /project/{project_id}/pod/{pod_id} | Update a pod&#39;s spec
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
    openapiclient "github.com/Podinate/podinate/lib/api_client"
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
    openapiclient "github.com/Podinate/podinate/lib/api_client"
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


## ProjectProjectIdPodPodIdExecPost

> string ProjectProjectIdPodPodIdExecPost(ctx, projectId, podId).Account(account).Command(command).Interactive(interactive).Tty(tty).Body(body).Execute()

Execute a command in a pod



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
    podId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    command := []string{"Inner_example"} // []string |  (optional)
    interactive := true // bool |  (optional) (default to false)
    tty := true // bool |  (optional) (default to false)
    body := os.NewFile(1234, "some_file") // *os.File |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodPodIdExecPost(context.Background(), projectId, podId).Account(account).Command(command).Interactive(interactive).Tty(tty).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPodIdExecPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodPodIdExecPost`: string
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodPodIdExecPost`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**podId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPodIdExecPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 
 **command** | **[]string** |  | 
 **interactive** | **bool** |  | [default to false]
 **tty** | **bool** |  | [default to false]
 **body** | ***os.File** |  | 

### Return type

**string**

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: application/octet-stream
- **Accept**: text/plain, application/json

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
    openapiclient "github.com/Podinate/podinate/lib/api_client"
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


## ProjectProjectIdPodPodIdLogsGet

> string ProjectProjectIdPodPodIdLogsGet(ctx, projectId, podId).Account(account).Lines(lines).Follow(follow).Execute()

Get the logs for a pod



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
    podId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    lines := int32(56) // int32 | The number of lines to return (optional) (default to 20)
    follow := true // bool | Whether to keep the connection open and continue streaming the logs (optional) (default to false)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodPodIdLogsGet(context.Background(), projectId, podId).Account(account).Lines(lines).Follow(follow).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPodIdLogsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodPodIdLogsGet`: string
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodPodIdLogsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**podId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPodIdLogsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **account** | **string** | The account to use for the request | 
 **lines** | **int32** | The number of lines to return | [default to 20]
 **follow** | **bool** | Whether to keep the connection open and continue streaming the logs | [default to false]

### Return type

**string**

### Authorization

[APIKeyAuth](../README.md#APIKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectProjectIdPodPodIdPut

> Pod ProjectProjectIdPodPodIdPut(ctx, projectId, podId).Account(account).Pod(pod).Execute()

Update a pod's spec



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
    podId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    pod := *openapiclient.NewPod("hello-world", "Hello World", "wordpress") // Pod | A JSON object containing the information needed to update a pod

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.PodApi.ProjectProjectIdPodPodIdPut(context.Background(), projectId, podId).Account(account).Pod(pod).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `PodApi.ProjectProjectIdPodPodIdPut``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectProjectIdPodPodIdPut`: Pod
    fmt.Fprintf(os.Stdout, "Response from `PodApi.ProjectProjectIdPodPodIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**projectId** | **string** |  | 
**podId** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectProjectIdPodPodIdPutRequest struct via the builder pattern


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
    openapiclient "github.com/Podinate/podinate/lib/api_client"
)

func main() {
    projectId := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    pod := *openapiclient.NewPod("hello-world", "Hello World", "wordpress") // Pod | A JSON object containing the information needed to create a new pod

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

