// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package cerr

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrBadPage             = Error("bad page")
	ErrBadReferrer         = Error("bad referrer")
	ErrCreateSchema        = Error("schema exists")
	ErrDatabaseExists      = Error("database exists")
	ErrForeignKeysDisabled = Error("foreign keys disabled")
	ErrMissingReferrer     = Error("missing referrer")
	ErrNotImplemented      = Error("not implemented")
	ErrPragmaReturnedNil   = Error("pragma returned nil")
	ErrUnknownLanguage     = Error("unknown language")
)
