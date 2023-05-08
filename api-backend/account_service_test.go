package main

import (
	"encoding/json"
	"math/rand"
	"testing"

	api "github.com/johncave/podinate/api-backend/go"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAccount(t *testing.T) {
	// TODO - implement TestRegisterAccount
	acc := api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	RegisterTestAccount(t, acc)
	DeleteTestAccount(t, acc.Id)
}

func TestRegisterTooLongAccount(t *testing.T) {
	// TODO - implement TestRegisterBadAccount
	acc := api.Account{Id: RandStringBytes(64), Name: RandStringBytes(128)}
	writer := makeRequest("POST", "/v0/account", acc, nil)
	assert.Equal(t, 400, writer.Code)
}

func TestRegisterSameAccountTwice(t *testing.T) {
	// TODO - implement TestRegisterSameAccountTwice
	// acc := api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	// RegisterTestAccount(t, acc)
	// writer := makeRequest("POST", "/v0/account", acc, nil)
	// assert.Equal(t, 400, writer.Code)

	// DeleteTestAccount(t, acc.Id)
}

func TestGetAccount(t *testing.T) {
	acc := api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	RegisterTestAccount(t, acc)
	responseAcc := GetTestAccount(t, acc.Id)
	assert.Equal(t, acc.Id, responseAcc.Id)
	assert.Equal(t, acc.Name, responseAcc.Name)
	DeleteTestAccount(t, acc.Id)
}

func BenchmarkGetAccount(b *testing.B) {
	acc := api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	_ = makeRequest("POST", "/v0/account", acc, nil)
	headers := map[string]string{"Account": acc.Id}
	for n := 0; n < b.N; n++ {
		makeRequest("GET", "/v0/account", nil, headers)
	}
}

func TestPatchAccount(t *testing.T) {
	acc := api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	RegisterTestAccount(t, acc)
	getAcc := GetTestAccount(t, acc.Id)
	assert.Equal(t, acc.Id, getAcc.Id)
	assert.Equal(t, acc.Name, getAcc.Name)
	// Update the account
	acc.Name = RandStringBytes(20)

	headers := map[string]string{"Account": acc.Id}
	writer := makeRequest("PATCH", "/v0/account", acc, headers)
	assert.Equal(t, 200, writer.Code)
	var patchAcc api.Account
	err := json.Unmarshal(writer.Body.Bytes(), &patchAcc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, acc.Id, patchAcc.Id)
	assert.Equal(t, acc.Name, patchAcc.Name)
	assert.NotEqual(t, getAcc.Name, patchAcc.Name)
	DeleteTestAccount(t, acc.Id)
}

// Helper functions for testing
func RegisterTestAccount(t *testing.T, acc api.Account) {
	writer := makeRequest("POST", "/v0/account", acc, nil)
	assert.Equal(t, 201, writer.Code)
	var respAcc api.Account
	err := json.Unmarshal(writer.Body.Bytes(), &respAcc)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, acc.Id, respAcc.Id)
	assert.Equal(t, acc.Name, respAcc.Name)
	t.Log("Registered account: ", respAcc)
}

// GetTestAccount - get an account from the database
func GetTestAccount(t *testing.T, id string) api.Account {
	headers := map[string]string{"Account": id}
	writer := makeRequest("GET", "/v0/account", nil, headers)
	assert.Equal(t, 200, writer.Code)
	var responseAcc api.Account
	err := json.Unmarshal(writer.Body.Bytes(), &responseAcc)
	if err != nil {
		t.Error(err)
	}
	t.Log("Got account: ", responseAcc)
	return responseAcc
}

func DeleteTestAccount(t *testing.T, id string) {
	headers := map[string]string{"Account": id}
	writer := makeRequest("DELETE", "/v0/account", nil, headers)
	assert.Equal(t, 202, writer.Code)
}

// Utility functions for test

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
