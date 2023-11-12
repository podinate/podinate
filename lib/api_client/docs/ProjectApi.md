# \ProjectApi

All URIs are relative to *https://api.podinate.com/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ProjectGet**](ProjectApi.md#ProjectGet) | **Get** /project | Returns a list of projects.
[**ProjectIdDelete**](ProjectApi.md#ProjectIdDelete) | **Delete** /project/{id} | Delete an existing project
[**ProjectIdGet**](ProjectApi.md#ProjectIdGet) | **Get** /project/{id} | Get an existing project given by ID
[**ProjectIdPatch**](ProjectApi.md#ProjectIdPatch) | **Patch** /project/{id} | Update an existing project
[**ProjectPost**](ProjectApi.md#ProjectPost) | **Post** /project | Create a new project



## ProjectGet

> []ProjectGet200ResponseInner ProjectGet(ctx).Account(account).Page(page).Limit(limit).Execute()

Returns a list of projects.



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
    account := "my-account" // string | The account to use for the request
    page := int32(56) // int32 | The page number to return (starts at 0) (optional) (default to 0)
    limit := int32(56) // int32 | The amount of items to return per page (optional) (default to 20)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectApi.ProjectGet(context.Background()).Account(account).Page(page).Limit(limit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectApi.ProjectGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectGet`: []ProjectGet200ResponseInner
    fmt.Fprintf(os.Stdout, "Response from `ProjectApi.ProjectGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProjectGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **account** | **string** | The account to use for the request | 
 **page** | **int32** | The page number to return (starts at 0) | [default to 0]
 **limit** | **int32** | The amount of items to return per page | [default to 20]

### Return type

[**[]ProjectGet200ResponseInner**](ProjectGet200ResponseInner.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectIdDelete

> Project ProjectIdDelete(ctx, id).Account(account).Execute()

Delete an existing project



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
    id := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectApi.ProjectIdDelete(context.Background(), id).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectApi.ProjectIdDelete``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectIdDelete`: Project
    fmt.Fprintf(os.Stdout, "Response from `ProjectApi.ProjectIdDelete`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 

### Return type

[**Project**](Project.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectIdGet

> Project ProjectIdGet(ctx, id).Account(account).Execute()

Get an existing project given by ID



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
    id := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectApi.ProjectIdGet(context.Background(), id).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectApi.ProjectIdGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectIdGet`: Project
    fmt.Fprintf(os.Stdout, "Response from `ProjectApi.ProjectIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 

### Return type

[**Project**](Project.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectIdPatch

> Project ProjectIdPatch(ctx, id).Account(account).Project(project).Execute()

Update an existing project



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
    id := "hello-world" // string | 
    account := "my-account" // string | The account to use for the request
    project := *openapiclient.NewProject() // Project |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectApi.ProjectIdPatch(context.Background(), id).Account(account).Project(project).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectApi.ProjectIdPatch``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectIdPatch`: Project
    fmt.Fprintf(os.Stdout, "Response from `ProjectApi.ProjectIdPatch`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiProjectIdPatchRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **account** | **string** | The account to use for the request | 
 **project** | [**Project**](Project.md) |  | 

### Return type

[**Project**](Project.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ProjectPost

> Project ProjectPost(ctx).Account(account).Project(project).Execute()

Create a new project



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
    account := "my-account" // string | The account to use for the request
    project := *openapiclient.NewProject() // Project | A JSON object containing the information needed to create a new project

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.ProjectApi.ProjectPost(context.Background()).Account(account).Project(project).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ProjectApi.ProjectPost``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ProjectPost`: Project
    fmt.Fprintf(os.Stdout, "Response from `ProjectApi.ProjectPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiProjectPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **account** | **string** | The account to use for the request | 
 **project** | [**Project**](Project.md) | A JSON object containing the information needed to create a new project | 

### Return type

[**Project**](Project.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

