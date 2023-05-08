package account

import (
	"errors"

	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

type Account struct {
	Uuid string
	Name string
	ID   string
}

// Validate account checks that a user's desired account properties are allowed
func (a *Account) ValidateNew() error {
	// check the account id and name are not too long
	if len(a.ID) > 50 {
		return errors.New("Account ID too long")
	}
	if len(a.Name) > 64 {
		return errors.New("Account name too long")
	}
	return nil
}

// Register creates a new account in the database
func (a *Account) Register() error {
	err := a.ValidateNew()
	if err != nil {
		return err
	}
	_, err = config.DB.Exec("INSERT INTO account(uuid, id, name) VALUES(gen_random_uuid(), $1, $2)", a.ID, a.Name)
	// Check if insert was successful
	if err != nil {
		return errors.New("Account ID not available")
	}
	return nil

}

// Patch updates an account in the database
func (a *Account) Patch(requested api.Account) error {
	// Check which fields are actually being updated
	if requested.Id != "" {
		a.ID = requested.Id
	}
	if requested.Name != "" {
		a.Name = requested.Name
	}
	// Update the database
	_, err := config.DB.Exec("UPDATE account SET name = $1, id = $2 WHERE uuid = $3", a.Name, a.ID, a.Uuid)
	if err != nil {
		return errors.New("Could not update account")
	}
	return nil
}

// GetBySlug retrieves an account from the database by its slug
func (a *Account) GetByID(desired_id string) error {
	row := config.DB.QueryRow("SELECT uuid, id, name FROM account WHERE id = $1 LIMIT 1", desired_id)

	err := row.Scan(&a.Uuid, &a.ID, &a.Name)
	if err != nil {
		return errors.New("Could not find this account")
	}
	return nil
}

func (a *Account) ToAPIAccount() api.Account {
	return api.Account{Id: a.ID, Name: a.Name}
}

// Delete removes an account from the database
func (a *Account) Delete() error {
	_, err := config.DB.Exec("DELETE FROM account WHERE uuid = $1", a.Uuid)
	if err != nil {
		return errors.New("Could not delete account")
	}
	return nil
}
