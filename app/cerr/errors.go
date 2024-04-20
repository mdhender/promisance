// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package cerr

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrCreateSchema        = Error("schema exists")
	ErrDatabaseExists      = Error("database exists")
	ErrForeignKeysDisabled = Error("foreign keys disabled")
	ErrPragmaReturnedNil   = Error("pragma returned nil")
)
