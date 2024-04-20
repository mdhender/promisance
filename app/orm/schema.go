// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package orm

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
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

// CreateSqliteDatabase creates the SQLite database and initializes it.
// It returns an error if the database already exists.
func CreateSqliteDatabase(dbName string) (*DB, error) {
	// if the database exist, return an error.
	if _, err := os.Stat(dbName); err == nil {
		return nil, fmt.Errorf("database exists")
	}
	dbSqlite, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
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

	// create the schema
	log.Printf("orm: creating database schema\n")
	if _, err = dbSqlite.Exec(ddlScript); err != nil {
		log.Printf("orm: failed to create database schema\n")
		return nil, errors.Join(cerr.ErrCreateSchema, err)
	}

	// todo: this should be the auto-migrate logic
	log.Printf("orm: todo: migrate database\n")

	log.Printf("orm: created %s", dbName)
	return db, nil
}

// OpenSqliteDatabase opens a SQLite database.
// It returns an error if the database does not exist.
func OpenSqliteDatabase(dbName string) (*DB, error) {
	// if the database does not exist, return an error.
	if _, err := os.Stat(dbName); err != nil {
		return nil, err
	}
	dbSqlite, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, err
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
	log.Printf("orm: todo: migrate database\n")

	log.Printf("orm: opened database %s", dbName)
	return db, nil
}
