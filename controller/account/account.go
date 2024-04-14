package account

import (
	"math/rand"
	"net/http"
	"regexp"
	"time"

	user "github.com/johncave/podinate/controller/user"

	"github.com/johncave/podinate/controller/apierror"
	"github.com/johncave/podinate/controller/config"
	api "github.com/johncave/podinate/controller/go"

	lh "github.com/johncave/podinate/controller/loghandler"
)

type Account struct {
	UUID string
	Name string
	ID   string
}

const (
	ActionView   = "account:view"
	ActionDelete = "account:delete"
	ActionUpdate = "account:update"
	ActionCreate = "account:create"
)

// Validate account checks that a user's desired account properties are allowed
func ValidateNewRequested(a api.Account) *apierror.ApiError {
	// check the account id and name are not too long
	if len(a.Id) > 30 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Account ID too long"}
	}
	if len(a.Name) > 64 {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Account name too long"}
	}

	lh.Log.Debugw("Validating account has correct ID", "account", a)
	m, err := regexp.MatchString(`^([a-z]*[0-9]*-*)*$`, a.Id)
	if err != nil {
		return &apierror.ApiError{Code: http.StatusBadRequest, Message: "Error checking project ID " + err.Error()}
	}
	if !m { // ID must be lowercase letters and numbers only
		return apierror.New(http.StatusBadRequest, "Account ID "+a.Id+" invalid, must be lowercase letters, numbers and dashes only")
	}
	lh.Log.Debugw("Validated account", "account", a)
	return nil
}

// CreateTest creates a new account in the database for testing purposes
func CreateTest() (*Account, *user.User, *apierror.ApiError) {
	rand.Seed(time.Now().UnixNano())

	// Generate random strings for test account ID and name
	id := generateRandomString(10)
	name := generateRandomString(10)

	req := api.Account{Id: id, Name: name}
	owner, _, err := user.Create(nil, nil, nil, nil)
	if err != nil {
		lh.Log.Errorw("Error getting test account owner", "error", err)
		return nil, nil, apierror.NewWithError(http.StatusInternalServerError, "Error getting test account owner ", err)
	}
	lh.Log.Debugw("Creating test account", "account", req, "owner", owner)
	a, err := Create(req, owner)
	lh.Log.Debug("Creating test account error", "error", err, "account", a)
	// I don't know why but this is inverted on my machine for some reason
	if err != nil {
		lh.Log.Errorw("Error creating test account", "error", err)
		return nil, nil, apierror.New(http.StatusInternalServerError, "Error creating test account "+err.Error())
	}
	return &a, owner, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Register creates a new account in the database
func Create(requestedAccount api.Account, owner *user.User) (Account, *apierror.ApiError) {
	lh.Log.Debugw("Creating account", "account", requestedAccount, "owner", owner)
	// TODO: Make this take an api.Account instead of an account.Account
	a := Account{ID: requestedAccount.Id, Name: requestedAccount.Name}

	err := ValidateNewRequested(requestedAccount)
	lh.Log.Debugw("Validation result is ", "account", requestedAccount, "error", err)
	if err != nil {
		lh.Log.Errorw("Error validating new account", "error", err, "account", a)
		return Account{}, err
		//apierror.New(http.StatusBadRequest, err.Error())
	} else {
		lh.Log.Debug("Validation passed", "account", a)
	}
	//_, dberr := config.DB.Exec("INSERT INTO account(uuid, id, name) VALUES(gen_random_uuid(), $1, $2)", a.ID, a.Name).Scan(&a.Uuid)
	dberr := config.DB.QueryRow("INSERT INTO account(uuid, id, name, owner_uuid) VALUES(gen_random_uuid(), $1, $2, $3) RETURNING uuid", a.ID, a.Name, owner.GetUUID()).Scan(&a.UUID)
	// Check if insert was successful
	if dberr != nil {
		lh.Log.Errorw("Error creating account", "error", dberr, "account", a)
		return Account{}, &apierror.ApiError{Code: http.StatusBadRequest, Message: dberr.Error()}
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
	_, err := config.DB.Exec("UPDATE account SET name = $1, id = $2 WHERE uuid = $3", a.Name, a.ID, a.UUID)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, "Could not update account")
	}
	return apierror.New(http.StatusOK, "")
}

// GetBySlug retrieves an account from the database by its slug
func GetByID(desired_id string) (Account, *apierror.ApiError) {
	row := config.DB.QueryRow("SELECT uuid, id, name FROM account WHERE id = $1 LIMIT 1", desired_id)

	var a Account
	err := row.Scan(&a.UUID, &a.ID, &a.Name)
	if err != nil {
		return Account{}, apierror.New(http.StatusNotFound, "Could not find account "+desired_id)
	}
	return a, nil
}

func (a *Account) ToAPIAccount() api.Account {
	return api.Account{Id: a.ID, Name: a.Name, ResourceId: a.GetResourceID()}
}

// Delete removes an account from the database
// Lol @ how complicated this function is gonna get
// Only if we code it poorly, John
func (a *Account) Delete() *apierror.ApiError {
	lh.Log.Debugw("Deleting account", "account", a)
	_, err := config.DB.Exec("DELETE FROM account WHERE uuid = $1", a.UUID)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, "Could not delete account "+err.Error())
	}
	return nil
}

// GetUUID returns the UUID of the account
func (a *Account) GetUUID() string {
	return a.UUID
}

// GetRID returns the RID of the account
func (a Account) GetResourceID() string {
	return "account:" + a.ID
}

// GetLabels returns the labels of the account
func (a *Account) GetLabels() map[string]string {
	return map[string]string{"podinate.com/account": a.ID}
}

// GetAnnotations returns the annotations of the account
func (a *Account) GetAnnotations() map[string]string {
	// This is blank for now, but could be very useful in the future
	return map[string]string{}
}
