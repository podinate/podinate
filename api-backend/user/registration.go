package user

import (
	"errors"

	"github.com/johncave/podinate/api-backend/config"
	"github.com/markbates/goth"

	lh "github.com/johncave/podinate/api-backend/loghandler"
)

func RegisterUser(user goth.User) (*User, error) {
	var authorisedUser string
	var providerID string

	// Check the user doesn't already exist
	err := config.DB.QueryRow("SELECT provider_id, authorised_user FROM oauth_login WHERE provider = $1 AND provider_id = $2", user.Provider, user.UserID).Scan(&providerID, &authorisedUser)
	if err != nil {
		lh.Log.Infow("User not found, inserting new user", "provider", user.Provider, "provider_id", user.UserID)
		// The user is not registered, so we need to register them
		// Add the user to the user table
		err = config.DB.QueryRow("INSERT INTO \"user\" (main_provider, id, display_name, email, avatar_url) VALUES ($1, $2, $3, $4, $5) RETURNING uuid", user.Provider, user.NickName, user.Name, user.Email, user.AvatarURL).Scan(&authorisedUser)
		if err != nil {
			lh.Log.Errorw("Error inserting new user into DB", "error", err)
			return &User{}, errors.New("Error inserting new user into DB")
		}

		// Insert into the oauth_login table
		_, err = config.DB.Exec("INSERT INTO oauth_login (provider, provider_id, provider_username, access_token, refresh_token, authorised_user) VALUES ($1, $2, $3, $4, $5, $6)", user.Provider, user.UserID, user.NickName, user.AccessToken, user.RefreshToken, authorisedUser)
		if err != nil {
			lh.Log.Errorw("Error inserting oauth_login", "error", err)
			return &User{}, errors.New("Error inserting oauth_login")
		}

		out, err := GetByUUID(authorisedUser)
		lh.Log.Infow("user registered", "user_uuid", authorisedUser, "user", out)
		return out, err
	}

	// The user is already registered, so we just need to update the oauth_login table
	_, err = config.DB.Exec("UPDATE oauth_login SET provider_username = $1, access_token = $2, refresh_token = $3 WHERE provider = $4 AND provider_id = $5", user.NickName, user.AccessToken, user.RefreshToken, user.Provider, user.UserID)
	if err != nil {
		lh.Log.Errorw("Error updating oauth_login", "error", err)
		return &User{}, errors.New("Error updating oauth_login")
	}

	// Update the user table
	_, err = config.DB.Exec("UPDATE \"user\" SET display_name = $1, email = $2, avatar_url = $3 WHERE uuid = $4", user.Name, user.Email, user.AvatarURL, authorisedUser)
	if err != nil {
		lh.Log.Errorw("Error updating user", "error", err, "user_uuid", authorisedUser)
		return &User{}, errors.New("Error updating user")
	}
	out, err := GetByUUID(authorisedUser)
	lh.Log.Infow("User updated", "reason", "new login, updated details from oauth", "user_uuid", authorisedUser, "user", user)
	return out, err
}
