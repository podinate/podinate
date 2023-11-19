# User

User package provides functions related to user accounts in the system. 

For now users are people logged in from our own Gitlab instance, given it's all Oauth, soon I'd like to add Github and Gitlab(proper) login. 

## Goth
We use a package called Goth to provide the Oauth heavy lifting. Goth is intentionally left out of this package, all of the code which interacts with Goth is in the Router package in the `user_auth_handlers.go` file. 

## Login
At a high level, logins work like this:
1) The client requests a login session. We respond with a token that identifies that session and a redirect URL.
2) The redirect url is /user/login/redirect/{token}. When they go to this URL a cookie is set with the session token and they are redirected to the provider.
3) The user authorizes the session with the provider and is redirected to our callback endpoint. Goth handles all of unwrapping the session and gives us a user. We set up the user in the database at this point, and the callback endpoint just displays some html telling the user to close their browser tab now.
4) The client should then call the /user/login/complete endpoint with the token. At that point they can exchange the token for a proper API key 