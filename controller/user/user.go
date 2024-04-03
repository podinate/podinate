package user

import (
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	api "github.com/johncave/podinate/controller/go"
	"github.com/matthewhartstonge/argon2"
)

// User is a user
type User struct {
	UUID         string
	MainProvider sql.NullString
	ID           string
	Email        string
	DisplayName  string
	AvatarURL    sql.NullString
	Flags        map[string]string
}

func GetByUUID(id string) (*User, error) {
	// Get from the database by UUID
	var userOut User
	err := config.DB.QueryRow("SELECT uuid, main_provider, id, email, display_name, avatar_url FROM \"user\" WHERE uuid = $1", id).Scan(&userOut.UUID, &userOut.MainProvider, &userOut.ID, &userOut.Email, &userOut.DisplayName, &userOut.AvatarURL)
	if err != nil {
		return nil, err
	}

	return &userOut, nil
}

func GetByUsername(username string) (*User, *apierror.ApiError) {
	// Get from the database by UUID
	var userOut User
	err := config.DB.QueryRow("SELECT uuid, main_provider, id, email, display_name, avatar_url FROM \"user\" WHERE id = $1", username).Scan(&userOut.UUID, &userOut.MainProvider, &userOut.ID, &userOut.Email, &userOut.DisplayName, &userOut.AvatarURL)
	if err != nil {
		return nil, apierror.NewWithError(http.StatusNotFound, "Could not find user", err)
	}

	return &userOut, nil
}

// CheckInternalLogin checks if a login with username and password for an internal user is correct
func CheckInternalLogin(user_id string, password string) (*User, error) {
	// Get from the database by UUID
	var userOut User
	var passwordHash string
	err := config.DB.QueryRow("SELECT uuid, main_provider, id, email, display_name, avatar_url, password_hash FROM \"user\" WHERE id = $1 AND main_provider IS NULL", user_id).Scan(&userOut.UUID, &userOut.MainProvider, &userOut.ID, &userOut.Email, &userOut.DisplayName, &userOut.AvatarURL, &passwordHash)
	if err != nil {
		return nil, err
	}

	// TODO - check password hash
	decode, err := base64.StdEncoding.DecodeString(passwordHash)
	if err != nil {
		return nil, err
	}

	ok, err := argon2.VerifyEncoded([]byte(password), decode)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apierror.New(403, "Invalid username or password")
	}

	return &userOut, nil
}

func (u *User) GetResourceID() string {
	if u.MainProvider.String == "" {
		return "user:" + u.ID
	}
	return "user:" + u.MainProvider.String + ":" + u.ID
}

func (u *User) GetUUID() string {
	return u.UUID
}

func (u *User) GetMainProvider() string {
	return u.MainProvider.String
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetDisplayName() string {
	return u.DisplayName
}

func (u *User) ToAPI() *api.User {
	return &api.User{
		ResourceId:  u.GetResourceID(),
		Email:       u.GetEmail(),
		DisplayName: u.GetDisplayName(),
		AvatarUrl:   u.AvatarURL.String,
	}
}
