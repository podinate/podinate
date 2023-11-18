package user

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/johncave/podinate/api-backend/config"
	eh "github.com/johncave/podinate/api-backend/errorhandler"
	"golang.org/x/crypto/sha3"
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

// IssueAPIKey issues an API key for the user.
// ClientName is the name of the client that the key is being issued for. For example "podinate-cli".
func (u *User) IssueAPIKey(clientName string) (string, error) {
	// Generate the key
	key, err := GenerateAPIKey()
	if err != nil {
		eh.Log.Errorw("Error generating API key", "error", err)
		return "", err
	}

	h := sha3.New512()
	h.Write([]byte(key))
	store := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Insert the key into the database
	_, err = config.DB.Exec("INSERT INTO api_key (user_uuid, name, key, expires) VALUES ($1, $2, $3, CURRENT_TIMESTAMP + interval '1 year')", u.UUID, clientName, store)
	if err != nil {
		eh.Log.Errorw("Error inserting API key", "error", err)
		return "", err
	}

	return "puak-" + key, nil
}

func GenerateAPIKey() (string, error) {
	// Generate 32 random characters
	buf := make([]byte, 64)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	eh.Log.Infow("Generated API key", "key", buf)

	// Base64 encode the random bytes
	b := base64.URLEncoding.EncodeToString(buf)
	return b, nil

}

func (u *User) GetRID() string {
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
