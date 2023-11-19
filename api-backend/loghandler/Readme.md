# loghandler
This package handles a bunch of stuff relating to logs, if something goes wrong, give the requestor a vague error message, and write as much context as possible to this package. It should be used like so,

## Use 
The primary use is the logging built on Uber's [zap](https://github.com/uber-go/zap) library.

Simply use as below:
```go
import lh "github.com/johncave/podinate/api-backend/loghandler"

lh.Debug(ctx, "error saving user",
    "user", user,
    "account", account,    
)
```
Give as much context as possible and prefer lower case keys. Using the methods provided directly by loghandler and passing a context will allow the package to insert a request ID for traceability. If you don't have access to the context, you can access the zap.SugaredLogger instance directly like so,
```go
import lh "github.com/johncave/podinate/api-backend/loghandler"

lh.Log.Debug("error saving user",
    "user", user,
    "account", account
)
```