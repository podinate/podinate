package account

import (
	"net/http"

	"github.com/johncave/podinate/api-backend/apierror"
	"github.com/johncave/podinate/api-backend/config"
	api "github.com/johncave/podinate/api-backend/go"
)

type Account struct {
	Uuid string
	Name string
	ID   string
}

// Validate account checks that a user's desired account properties are allowed
func (a *Account) ValidateNew() *apierror.ApiError {
	// check the account id and name are not too long
	if len(a.ID) > 30 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Account ID too long"}
	}
	if len(a.Name) > 64 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Account name too long"}
	}
	return nil
}

// Register creates a new account in the database
func Create(requestedAccount api.Account) (Account, *apierror.ApiError) {
	// TODO: Make this take an api.Account instead of an account.Account
	a := Account{ID: requestedAccount.Id, Name: requestedAccount.Name}

	err := a.ValidateNew()
	if err != nil {
		return Account{}, err
		//apierror.New(http.StatusBadRequest, err.Error())
	}
	//_, dberr := config.DB.Exec("INSERT INTO account(uuid, id, name) VALUES(gen_random_uuid(), $1, $2)", a.ID, a.Name).Scan(&a.Uuid)
	dberr := config.DB.QueryRow("INSERT INTO account(uuid, id, name) VALUES(gen_random_uuid(), $1, $2) RETURNING uuid", a.ID, a.Name).Scan(&a.Uuid)

	// Check if insert was successful
	if dberr != nil {
		return Account{}, &apierror.ApiError{Code: http.StatusBadRequest, Message: err.Error()}

	}
	return a, nil

}

// Patch updates an account in the database
func (a *Account) Patch(requested api.Account) *apierror.ApiError {
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
		return apierror.New(http.StatusInternalServerError, "Could not update account")
	}
	return apierror.New(http.StatusOK, "")
}

// GetBySlug retrieves an account from the database by its slug
func GetByID(desired_id string) (Account, *apierror.ApiError) {
	row := config.DB.QueryRow("SELECT uuid, id, name FROM account WHERE id = $1 LIMIT 1", desired_id)

	var a Account
	err := row.Scan(&a.Uuid, &a.ID, &a.Name)
	if err != nil {
		return Account{}, apierror.New(http.StatusNotFound, "Could not find account")
	}
	return a, nil
}

func (a *Account) ToAPIAccount() api.Account {
	return api.Account{Id: a.ID, Name: a.Name}
}

// Delete removes an account from the database
// Lol @ how complicated this function is gonna get
func (a *Account) Delete() *apierror.ApiError {
	_, err := config.DB.Exec("DELETE FROM account WHERE uuid = $1", a.Uuid)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, "Could not delete account")
	}
	return nil
}
