# Router
This is a router which implements all of the servicer functions from the api package, which has the Service methods. 

## Service methods 
The API package in `../go` has various "*APIServices", these are basically router handlers, which handle unwrapping parameters from the API request and passing them into an "*APIServicer". 

## Servicer
We create a custom "*APIServicer" to pass to the generated API package's "*APIService", for example look in `account_handlers.go`. These functions handle the job of retrieving the structures in question and calling the right methods, and also wrapping up any errors into tidy http errors. Any validation should *not* be done here, it should be done in the package that handles that struct. For example, `AccountGetHandler()` should call `account.GetByID()` and pass any errors to the client, then return `a.ToAPI()` and set the error code to 200.  

## Auth
Auth should **not** be handled anywhere in this package, for information about how auth works, check `iam/Readme.md`.