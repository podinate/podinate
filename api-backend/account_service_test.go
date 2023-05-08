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

}

func TestGetAccount(t *testing.T) {
	// TODO - implement TestGetAccount
	acc := api.Account{Id: RandStringBytes(10), Name: RandStringBytes(20)}
	RegisterTestAccount(t, acc)
	writer := makeRequest("GET", "/v0/account/"+acc.Id, nil)
	assert.Equal(t, 200, writer.Code)
	var responseAcc api.Account
	json.Unmarshal(writer.Body.Bytes(), &responseAcc)
	assert.Equal(t, acc.Id, responseAcc.Id)
	assert.Equal(t, acc.Name, responseAcc.Name)
}

func TestPatchAccount(t *testing.T) {
	// TODO - implement TestPatchAccount
	t.Skip("Not implemented")
}

func RegisterTestAccount(t *testing.T, acc api.Account) {
	writer := makeRequest("POST", "/v0/account", acc)
	assert.Equal(t, 201, writer.Code)
	var respAcc api.Account
	json.Unmarshal(writer.Body.Bytes(), &respAcc)
	assert.Equal(t, acc.Id, respAcc.Id)
	assert.Equal(t, acc.Name, respAcc.Name)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
