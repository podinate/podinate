package account

import (
	"os"
	"regexp"
	"testing"

	api "github.com/johncave/podinate/controller/go"
)

func TestValidation(t *testing.T) {
	new := api.Account{Id: "test", Name: "Test Account"}
	err := ValidateNewRequested(new)
	if err != nil {
		t.Errorf("Validation failed: %s", err)
	}
}

func TestValidationTooLongID(t *testing.T) {
	new := api.Account{Id: "test123456789012345678901234567890", Name: "Test Account"}
	err := ValidateNewRequested(new)
	if err == nil {
		t.Errorf("Validation should have failed for too long account ID")
	}
}

func TestValidationTooLongName(t *testing.T) {
	new := api.Account{Id: "test", Name: "Test Account123456789012345678901234567890123456789012345678901234567890 lorem ipum lorem ipsum lorem ipsum"}
	err := ValidateNewRequested(new)
	if err == nil {
		t.Errorf("Validation should have failed for too long account name")
	}
}

func TestValidationInvalidID(t *testing.T) {
	new := api.Account{Id: "TestAccountWithCapitals", Name: "Test Account"}
	err := ValidateNewRequested(new)
	if err == nil {
		t.Errorf("Validation should have failed for invalid account ID")
	}
}

func TestValidationRegex(t *testing.T) {
	r := regexp.MustCompile(`^([a-z]*[0-9]*-*)*$`)
	if !r.MatchString("test") {
		t.Errorf("Regex should have matched")
	}
	if r.MatchString("TestAccountWithCapitals") {
		t.Errorf("Regex should not have matched")
	}

}

// func TestCreateAccountWithInvalidID(t *testing.T) {
// 	new := api.Account{Id: "TestAccountWithCapitals", Name: "Test Account"}
// 	owner, err := user.GetByUsername("administrator")
// 	if err != nil {
// 		t.Errorf("Error getting user: %s", err)
// 	}
// 	created, err := Create(new, owner)
// 	if err == nil {
// 		t.Errorf("Create should have failed for invalid account ID")
// 	}
// 	err = created.Delete()
// 	if err != nil {
// 		t.Errorf("Error deleting account: %s", err)
// 	}
// }

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
