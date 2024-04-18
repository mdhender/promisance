// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import (
	"log"
)

func (p *PHP) includes_PasswordHash_php() error {
	if !p.globals.IN_GAME {
		p.die("Access denied")
	}

	log.Printf("todo: includes/PasswordHash should add an instance to the engine\n")

	return nil
}

type PasswordHasher interface {
	PasswordHash(iterationCountLog2 int, portableHashes int)
	GetRandomBytes(count int) ([]byte, error)
	Encode64(input string, count int) (string, error)
	GenSaltPrivate(input string) (string, error)
	CryptPrivate(password, setting string) (string, error)
	GenSaltExtended(input string) (string, error)
	GenSaltBlowfish(input string) (string, error)
	HashPassword(password string) (string, error)
	CheckPassword(password string, hash string) (bool, error)
}
