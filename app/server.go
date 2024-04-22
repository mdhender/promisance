// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm"
)

type server struct {
	data          string // path to store data files
	templates     string // path to template files
	addr          string
	host          string
	port          string
	tz            string
	db            *orm.DB
	world         *model.World_t
	authenticator *AuthenticationManager
	sessions      *SessionManager_t
}

type AuthenticationManager struct{}
type User struct{}

func (a *AuthenticationManager) Authenticate(username, password string) (User, error) {
	return User{}, nil
}
