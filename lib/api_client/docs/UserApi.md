# \UserApi

All URIs are relative to *https://api.podinate.com/v0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UserGet**](UserApi.md#UserGet) | **Get** /user | Get the current user
[**UserLoginCallbackProviderGet**](UserApi.md#UserLoginCallbackProviderGet) | **Get** /user/login/callback/{provider} | User login callback URL for oauth providers
[**UserLoginCompleteGet**](UserApi.md#UserLoginCompleteGet) | **Get** /user/login/complete | Complete a user login
[**UserLoginInitiateGet**](UserApi.md#UserLoginInitiateGet) | **Get** /user/login/initiate | Get a login URL
[**UserLoginRedirectTokenGet**](UserApi.md#UserLoginRedirectTokenGet) | **Get** /user/login/redirect/{token} | User login redirect URL to oauth providers



## UserGet

> UserGet200Response UserGet(ctx).Account(account).Execute()

Get the current user



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

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserApi.UserGet(context.Background()).Account(account).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserApi.UserGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserGet`: UserGet200Response
    fmt.Fprintf(os.Stdout, "Response from `UserApi.UserGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **account** | **string** | The account to use for the request | 

### Return type

[**UserGet200Response**](UserGet200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLoginCallbackProviderGet

> string UserLoginCallbackProviderGet(ctx, provider).Execute()

User login callback URL for oauth providers



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
    provider := "github" // string | The oauth provider to use. Valid options will be github / gitlab / podinate, during alpha only podinate is allowed.

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserApi.UserLoginCallbackProviderGet(context.Background(), provider).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserApi.UserLoginCallbackProviderGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLoginCallbackProviderGet`: string
    fmt.Fprintf(os.Stdout, "Response from `UserApi.UserLoginCallbackProviderGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**provider** | **string** | The oauth provider to use. Valid options will be github / gitlab / podinate, during alpha only podinate is allowed. | 

### Other Parameters

Other parameters are passed through a pointer to a apiUserLoginCallbackProviderGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/html, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLoginCompleteGet

> UserLoginCompleteGet200Response UserLoginCompleteGet(ctx).Token(token).Execute()

Complete a user login



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
    token := "abc1234" // string | The token given by /user/login/init to get the user's actual API key once they have completed the oauth flow (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserApi.UserLoginCompleteGet(context.Background()).Token(token).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserApi.UserLoginCompleteGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLoginCompleteGet`: UserLoginCompleteGet200Response
    fmt.Fprintf(os.Stdout, "Response from `UserApi.UserLoginCompleteGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLoginCompleteGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **token** | **string** | The token given by /user/login/init to get the user&#39;s actual API key once they have completed the oauth flow | 

### Return type

[**UserLoginCompleteGet200Response**](UserLoginCompleteGet200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLoginInitiateGet

> UserLoginInitiateGet200Response UserLoginInitiateGet(ctx).Provider(provider).Execute()

Get a login URL



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
    provider := "github" // string | The oauth provider to use. Valid options will be github / gitlab / podinate, during alpha only podinate is allowed. (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserApi.UserLoginInitiateGet(context.Background()).Provider(provider).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserApi.UserLoginInitiateGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLoginInitiateGet`: UserLoginInitiateGet200Response
    fmt.Fprintf(os.Stdout, "Response from `UserApi.UserLoginInitiateGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUserLoginInitiateGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **string** | The oauth provider to use. Valid options will be github / gitlab / podinate, during alpha only podinate is allowed. | 

### Return type

[**UserLoginInitiateGet200Response**](UserLoginInitiateGet200Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UserLoginRedirectTokenGet

> string UserLoginRedirectTokenGet(ctx, token).Execute()

User login redirect URL to oauth providers



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
    token := "abc1234" // string | The token given by /user/login/init to get the user's actual API key once they have completed the oauth flow

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.UserApi.UserLoginRedirectTokenGet(context.Background(), token).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `UserApi.UserLoginRedirectTokenGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UserLoginRedirectTokenGet`: string
    fmt.Fprintf(os.Stdout, "Response from `UserApi.UserLoginRedirectTokenGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**token** | **string** | The token given by /user/login/init to get the user&#39;s actual API key once they have completed the oauth flow | 

### Other Parameters

Other parameters are passed through a pointer to a apiUserLoginRedirectTokenGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/html, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
