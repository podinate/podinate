# api-backend
This is the backend for the API, in charge of serving the API. Please take a moment to read this tour. 

## Tour of the Code
#### Main
The main package is mostly empty, it sets up a router from the `router` package, then listenandserves(). 
#### Router 
Router package sets up a bunch of different (ApiService), for example UserApiService which are structs built by the OpenAPI code generator. For the most part we use the default service, as it will extract and provide paramaters easily to our Servicer methods, and provide our own ApiServicer, for example `UserApiServicer`. Have a look in the `go/` folder to see these services.
#### Logging 
Logging is handled by the `loghandler` package. It is built on the [zap](https://github.com/uber-go/zap) library from Uber. You can access the zap sugared logger interface like so,
```go
import lh "github.com/johncave/podinate/api-backend/loghandler"
lh.Info(ctx, "error creating user", 
    "user", user, 
    "account", account
)
```
You can also access the raw zap.SugaredLogger instance at `lh.Log`

Where possible, please log the request ID. 
#### Authorization
Authorization is all handled by the `iam` package. You must read the readme of the `iam` package. 
#### Authentication
Currently handled with API keys, have a look at `users` package to see some of how that works. 
