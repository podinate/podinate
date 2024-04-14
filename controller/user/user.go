package user

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/google/uuid"
	"github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	api "github.com/johncave/podinate/controller/go"
	"github.com/matthewhartstonge/argon2"
	"go.tmthrgd.dev/passit"
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

// Create creates a new user in the database
func Create(username *string, password *string, email *string, displayName *string) (*User, *string, *apierror.ApiError) {
	// Generate a random UUID
	uuid := uuid.New().String()
	var err error

	if username == nil {
		// Generate a username
		u, err := passit.Repeat(passit.EFFLargeWordlist, "-", 2).Password(rand.Reader)
		if err != nil {
			return nil, nil, apierror.NewWithError(http.StatusInternalServerError, "Error generating username", err)
		}
		code, err := passit.Repeat(passit.Digit, "", 4).Password(rand.Reader)
		if err != nil {
			return nil, nil, apierror.NewWithError(http.StatusInternalServerError, "Error generating username", err)
		}
		u = u + code
		username = &u
	}

	if password == nil {
		// Generate a password
		p, err := passit.Repeat(passit.EFFLargeWordlist, " ", 5).Password(rand.Reader)
		if err != nil {
			return nil, nil, apierror.NewWithError(http.StatusInternalServerError, "Error generating password", err)
		}
		password = &p
	}

	if displayName == nil {
		displayName = username
	}

	if email == nil {
		e := "test@podinate.com"
		email = &e
	}

	// Hash the password
	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(*password))
	if err != nil {
		return nil, nil, apierror.NewWithError(http.StatusInternalServerError, "Error hashing password", err)
	}
	store := base64.StdEncoding.EncodeToString(hash)

	// Insert the user into the database
	var userUUID string
	err = config.DB.QueryRow("INSERT INTO \"user\" (uuid, id, email, display_name, password_hash) VALUES ($1, $2, $3, $4, $5) RETURNING uuid", uuid, *username, *email, *displayName, store).Scan(&userUUID)
	if err != nil {
		return nil, nil, apierror.NewWithError(http.StatusInternalServerError, "Error creating user", err)
	}

	// Return the user
	user := User{
		UUID:        uuid,
		ID:          *username,
		Email:       *email,
		DisplayName: *displayName,
	}
	return &user, password, nil
}

// Delete deletes a user from the database
func (u *User) Delete(ctx context.Context) *apierror.ApiError {
	_, err := config.DB.Exec("DELETE FROM \"user\" WHERE uuid = $1", u.UUID)
	if err != nil {
		return apierror.NewWithError(http.StatusInternalServerError, "Error deleting user", err)
	}

	return nil
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
