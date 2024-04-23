// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package authn implements a primitive and naive authentication package.
package authn

import (
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm"
	"log"
)

type Authenticator struct {
	db *orm.DB
}

func New(db *orm.DB) (*Authenticator, error) {
	a := &Authenticator{
		db: db,
	}
	return a, nil
}

func (a *Authenticator) Authenticate(username string, password string) (*model.User_t, error) {
	// todo: use bcrypt to encrypt the password
	log.Printf("authn: please use bcrypt to encrypt the password!\n")
	user, err := a.db.AuthenticatedUserFetch(username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *Authenticator) UserRoles(user *model.User_t) map[string]bool {
	roles := map[string]bool{}
	if user == nil {
		return roles
	}
	if user.Flags.Admin {
		roles["admin"] = true
	}
	if user.Flags.Closed {
		roles["closed"] = true
	}
	if user.Flags.Disabled {
		roles["disabled"] = true
	}
	if user.Flags.Mod {
		roles["mod"] = true
	}
	if user.Flags.Valid {
		roles["valid"] = true
	}
	if user.Flags.Watch {
		roles["watch"] = true
	}
	return roles
}
