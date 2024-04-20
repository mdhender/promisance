// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package orm

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"github.com/mdhender/promisance/app/cerr"
	"github.com/mdhender/promisance/app/orm/sqlc"
	"github.com/mdhender/semver"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var (
	// schema version
	schemaVersion = semver.Version{Major: 0, Minor: 1}
	//go:embed sqlc/sqlite/ddl.sql
	ddlScript string
)

type DB struct {
	ctx      context.Context
	db       *sqlc.Queries
	dbSqlite *sql.DB
}

func (db *DB) Close() error {
	var err error
	if db.dbSqlite != nil {
		err = db.dbSqlite.Close()
		db.dbSqlite = nil
	}
	return err
}

// OpenSqliteDatabase opens a SQLite database.
// If the database does not exist, it creates it and
// runs the data initialization scripts (in other words,
// it creates and then runs the migration scripts).
func OpenSqliteDatabase(dbName string) (*DB, error) {
	// if the database does not exist, create it and migrate it.
	dbCreate, dbMigrate := false, false
	if _, err := os.Stat(dbName); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		dbCreate, dbMigrate = true, true
	}
	dbSqlite, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
	}
	if dbCreate {
		log.Printf("orm: created %s", dbName)
	} else {
		log.Printf("orm: opened %s", dbName)
	}
	db := &DB{
		ctx:      context.Background(),
		db:       sqlc.New(dbSqlite),
		dbSqlite: dbSqlite,
	}

	// confirm that the database has foreign keys enabled
	var rslt sql.Result
	checkPragma := "PRAGMA" + " foreign_keys = ON"
	if rslt, err = dbSqlite.Exec(checkPragma); err != nil {
		log.Printf("orm: error: foreign keys are disabled\n")
		return nil, cerr.ErrForeignKeysDisabled
	} else if rslt == nil {
		log.Printf("orm: error: foreign keys pragma failed\n")
		return nil, cerr.ErrPragmaReturnedNil
	}

	// todo: this should be the auto-migrate logic
	// create the schema if needed
	if dbCreate {
		log.Printf("orm: initializing database\n")
		if _, err = dbSqlite.Exec(ddlScript); err != nil {
			log.Printf("orm: failed to initialize schema\n")
			return nil, errors.Join(cerr.ErrCreateSchema, err)
		}
	}
	if dbMigrate {
		log.Printf("orm: todo: migrate database\n")
	}

	return db, nil
}
