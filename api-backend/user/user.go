package user

import (
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

// User is a user
type User struct {
	UUID         string
	MainProvider string
	ID           string
	Email        string
	DisplayName  string
	AvatarURL    string
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

func (u *User) GetResourceID() string {
	return "user:" + u.MainProvider + ":" + u.ID
}

func (u *User) GetUUID() string {
	return u.UUID
}

func (u *User) GetMainProvider() string {
	return u.MainProvider
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
		AvatarUrl:   u.AvatarURL,
	}
}
