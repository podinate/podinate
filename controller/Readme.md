# controller
This is the backend for the API, in charge of serving the API. Please take a moment to read this tour. 

## Tour of the Code
#### Main
The main package is nearly empty, it just calls the cmd package. 
#### cmd
The cmd package handles command line options for the controller. Currently there are only two commands: 
- **controller** (no options) - Run the Podinate controller and listen for requests
- **controller init --email --ip** - Initilise the default user and account, write the profile information to /profile.yaml. 
#### iam
You must read the readme of the `iam` package. 
The IAM package is a barebones implementation of a policy-based IAM system.
Authorization is all handled by the `iam` package. 
#### go
This is the generated OpenAPI server code. You should directly use the Api**Service** structs from this package, and implement your own versions of the Api**Servicer** interfaces that they expect. This means you don't have to manually write logic to extract values from the request. If you need more control or direct access to the http.ResponseWriter or http.Request object, take a look at the `router/router.go` to see how to override some of the routes. 
#### Router 
The Router package has its own Readme that I recommend having a read of. 
#### Logging 
Logging is handled by the `loghandler` package. It is built on the [zap](https://github.com/uber-go/zap) library from Uber. You can access the zap sugared logger interface like so,
```go
import lh "github.com/johncave/podinate/controller/loghandler"
lh.Info(ctx, "error creating user", 
    "user", user, 
    "account", account
)
```
By passing the request context, the LogHandler package will grab out a request ID and throw it in the logged line. You can also access the raw zap.SugaredLogger instance at `lh.Log`. If logging manually, where possible, please log the request ID. 

#### Authentication
Currently handled with API keys, have a look at `users` package to see some of how that works. 
