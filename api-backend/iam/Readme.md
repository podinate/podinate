# IAM
I've tried to make this IAM package as easy and extensible as possible. 

## Overview 
The IAM system is a tiny implementation of Google Zanzibar reseach paper, conceptually similar to AWS IAM. 

#### Resource
A resource is anything in the system, be it a user, a pod, an account. 
#### Action
An action is something that a requestor can do. It should loosely correspond with an API endpoint. For example "GET /account" would correspond with "account:view".
#### Resource ID
Every Resource has a unique string identifier, an example of which would be "account:mine/project:my-blog/pod:frontend". Every potential resource must implement the iam.Resource interface, which just means having a `GetRID()` function. 
#### Requestor
A requestor is a resource which is requesting to do something. This is usually a user, but could be for example a pod running within the system. 
#### User 
Users are just another resource to tinyIAM. Everything specific to accounts **must** be in the user package instead. 

## How To Use
The iam package is plugged into the Router as a middleware. It gets the requesting resource by being passed the authorization header, then this resource is addded to the context of the request. To see if a requestor can do something, simply call `iam.RequestorCan(context, account.Account, resource, action)`. For an example of this, check `AccountGet()` in `router/account_handlers.go`. 

## Standards 
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