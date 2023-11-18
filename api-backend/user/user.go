package user

import "github.com/johncave/podinate/api-backend/config"

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
	err := config.DB.QueryRow("SELECT uuid, main_provider, id, email, display_name, avatar_url, flags FROM users WHERE uuid = $1", id).Scan(&userOut.UUID, &userOut.MainProvider, &userOut.ID, &userOut.Email, &userOut.DisplayName, &userOut.AvatarURL, &userOut.Flags)
	if err != nil {
		return nil, err
	}

	return &userOut, nil
}

func (u *User) GetRID() string {
	return ""
}
