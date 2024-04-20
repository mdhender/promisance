// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm"
)

type server struct {
	data  string // path to store data files
	addr  string
	host  string
	port  string
	tz    string
	db    *orm.DB
	world *model.World_t
}
