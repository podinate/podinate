package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/johncave/podinate/controller/config"
	"golang.org/x/crypto/sha3"

	lh "github.com/johncave/podinate/controller/loghandler"
)

// IssueAPIKey issues an API key for the user.
// ClientName is the name of the client that the key is being issued for. For example "podinate-cli".
func (u *User) IssueAPIKey(clientName string) (string, error) {
	// Generate the key
	key, err := GenerateAPIKey()
	if err != nil {
		lh.Log.Errorw("Error generating API key", "error", err)
		return "", err
	}

	h := sha3.New512()
	h.Write([]byte(key))
	store := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Insert the key into the database
	_, err = config.DB.Exec("INSERT INTO api_key (user_uuid, name, key, expires) VALUES ($1, $2, $3, CURRENT_TIMESTAMP + interval '1 year')", u.UUID, clientName, store)
	if err != nil {
		lh.Log.Errorw("Error inserting API key", "error", err)
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
	lh.Log.Infow("Generated API key", "key", buf)

	// Base64 encode the random bytes
	b := base64.URLEncoding.EncodeToString(buf)
	return b, nil

}

// GetFromAPIKey gets a user from an API key
func GetFromAPIKey(key string) (*User, error) {
	if len(key) < 5 {
		return nil, errors.New("invalid API key")
	}

	// Check the prefix
	if key[0:5] != "puak-" {
		return nil, errors.New("invalid API key")
	}

	// Remove the "puak-" prefix
	key = key[5:]

	// Hash the key
	h := sha3.New512()
	h.Write([]byte(key))
	store := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Get the key from the database
	var userUUID string
	err := config.DB.QueryRow("SELECT user_uuid FROM api_key WHERE key = $1", store).Scan(&userUUID)
	if err != nil {
		return nil, err
	}

	// Get the user from the database
	user, err := GetByUUID(userUUID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
