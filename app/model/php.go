// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package model

import "time"

type User_t struct {
	Id         int
	UserName   string
	Password   string
	Flags      UserFlag_t
	Nickname   string
	Email      string
	Lang       string
	DateFormat string
	LastIP     string
	CreateDate time.Time
	LastDate   time.Time
}

type UserFlag_t struct {
	// user has Moderator privileges (can set/clear multi and disabled flags, can browse empire messages)
	Mod bool
	// user has Administrator privileges (can grant/revoke privileges, delete/rename empires, login as anyone, edit clans)
	Admin bool
	// user account is disabled, cannot create new empires (but can still login to existing ones)
	Disabled bool
	// user account's email address has been validated at least once
	Valid bool
	// user account has been voluntarily closed, cannot create new empires or login to existing ones
	Closed bool
	// user account is suspected of abuse
	Watch bool
}
