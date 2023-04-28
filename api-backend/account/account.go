package account

import (
	"errors"

	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

type Account struct {
	Id   string
	Name string
	Slug string
}

// Validate account checks that a user's desired account properties are allowed
func (a *Account) ValidateNew() error {
	// check the account id and name are not too long
	if len(a.Slug) > 50 {
		return errors.New("Account ID too long")
	}
	if len(a.Name) > 64 {
		return errors.New("Account name too long")
	}
	return nil
}

func (a *Account) GetBySlug(slug string) error {
	row := config.DB.QueryRow("SELECT id, slug, name FROM account WHERE slug = $1", slug)
	var id, name string
	err := row.Scan(&id, &slug, &name)
	if err != nil {
		return errors.New("Could not find this account")
	}
	a.Id = id
	a.Name = name
	a.Slug = slug
	return nil
}

func (a *Account) ToAPIAccount() api.Account {
	return api.Account{Id: a.Slug, Name: a.Name}
}
