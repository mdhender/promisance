// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/promisance/app/jot"
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm"
)

type server struct {
	data          string // path to store data files
	templates     string // path to template files
	public        string // path to public files (sometimes called "assets")
	addr          string
	host          string
	port          string
	tz            string
	db            *orm.DB
	world         *model.World_t
	sessions      *jot.Factory_t
	authenticator *Authenticator_t
}

type Authenticator_t struct {
	user string
	pass string
}

func (a *Authenticator_t) Authenticate(username, password string) bool {
	return username == a.user && password == a.pass
}
