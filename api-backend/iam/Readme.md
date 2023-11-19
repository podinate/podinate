# IAM
I've tried to make this IAM package as easy and extensible as possible. 

## Overview 
The IAM system is a tiny implementation of Google Zanzibar reseach paper, conceptually similar to AWS IAM. 

#### Resource
A resource is anything in the system, be it a user, a pod, an account. 
#### Action
An action is something that a requestor can do. It should loosely correspond with an API endpoint. For example "GET /account" would correspond with "account:view".
#### Policy
A policy is a document stored in the database, each stores multiple versions so the administrator can easily roll back or forward to test potentially breaking changes.
#### Policy Document
A policy document is a YAML string which has one or more statements which allow or deny a given action, based on the requested Resource ID and Action. See the [Policy Document Format](#Policy Document Format) section.
#### Resource ID
Every Resource has a unique string identifier, an example of which would be "account:mine/project:my-blog/pod:frontend". Every potential resource must implement the iam.Resource interface, which just means having a `GetResourceID()` function which returns a resource ID string in this format.
#### Requestor
A requestor is a resource which is requesting to do something. This is usually a user, but could be for example a pod running within the system. 
#### User 
Users are just another resource to tinyIAM. Everything specific to humans logging in with Github/Gitlab **must** be in the user package instead. 

## How to Use 
Your handler function will not be called unless there's a valid requestor (unless it's an anonymous endpoint). Load up the resource the call would like to operate on then to see if a requestor can do something, simply call `iam.RequestorCan(context, account.Account, resource, action) bool`. For an example of this, check `AccountGet()` in `router/account_handlers.go`. `iam.RequestorCan` only returns boolean, if it encounters any error the result is false and the error is logged. 

## Your Code 
A resource should have several actions on it, these should be stored as constants in the corresponding package (don't hard code them bro). For example: 
```go
package database

const (
	ActionView = "database:view"
	ActionDelete = "database:delete"
	ActionUpdate = "database:update"
	ActionCreate = "database:create"
    ActionEnterCLI = "database:getcli"
)
```

## How It Works
Here's how a request is authorized end to end;
### Authentication
1) The router has an auth middleware. This middleware first checks if the user is calling one of the anonymous endpoints and allows it if so. 
2) The router then calls the iam package on `iam.AddRequestorToRequest(r *http.Request)`. This func knows how to authenticate a resource based on the request. If it can't add a requestor, the request is denied. 
3) Currently the only possible type of requestor is a user, it finds the user's API key, and loads up the user as the requestor.. In the future this is the method that would need to be expanded to allow other resources to authenticate. 
4) This method inserts the requestor into the context of the request, which is passed down to our handler methods. 
5) Your handler method, for example `AccountGet(ctx context.Context, accountID string)` will not be called unless a valid requestor is present.

### Authorization
1) Your handler method loads up the resource that it wants to operate on, and the account it's in, and calls`iam.RequestorCan(context.Context, account.Account, resource, action) bool`.
2) iam unpacks the requestor from the given context and looks for any policies the given account has assigned to the requestor.
3) For every attached policy, the policy document is unwrapped and compared against the requested resource's ID and the action being performed. The string matching is perfomed using the [github.com/gobwas/glob](https://github.com/gobwas/glob) package. 
4) If any policy statement matches, return true, if none do, return false. Repeat this for every policy. 
5) Return true or false. If any error is encountered we return false and log the error. 

## Policy Document Format
A policy document is a YAML file in the following form,
```yaml
version: 2023.1
statements:
  - effect: allow
    actions: ["**"]
    resources: ["**"]
```
Version is a field to allow future proofing. If we want to implement a new policy format in the future, we will create a new version. For now this must be `2023.1`. 

Statements is an array of statements that is evaluated in order they appear. They have the following fields:  
1) `effect` which is either `allow` or `deny`. Anything other than 'allow' is a denial. 
2) `actions` is an array of actions the statement applies to.
3) `resources` is an array of resource ID globs that the statement applies to. 

Let's consider the following more advanced example of a policy document `account:mine/project:my-blog/pod:the-blog` and action `pod:delete`.
```yaml
version: 2023.1
statements:
  - effect: deny
	actions:
	  - pod:delete
	resource:
	  - account:mine/project:my-blog/pod:*
  - effect: allow
	actions:
	  - **
	resource:
	  - account:mine/project:my-blog/**
```

The user has full access to all of the resources in the project `account:mine/project:my-blog`, but cannot delete any pods created in that project. 

Check the documentation for the [github.com/gobwas/glob](https://github.com/gobwas/glob) package for more information about how the globbing works.