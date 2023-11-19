package iam

import (
	user "github.com/johncave/podinate/api-backend/user"

	"github.com/johncave/podinate/api-backend/account"
)

func CanDo(account account.Account, user user.User, resource resource, action string) bool {
	return true
}

type resource interface {
	GetRID() string
}
