// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: dml.sql

package sqlc

import (
	"context"
	"database/sql"
)

const clanFetch = `-- name: ClanFetch :one

SELECT c_id,
       c_name,
       c_password,
       c_members,
       e_id_leader,
       e_id_asst,
       e_id_fa1,
       e_id_fa2,
       c_title,
       c_url,
       c_pic
FROM clan
WHERE c_id = ?
`

// Copyright (c) 2024 Michael D Henderson. All rights reserved.
func (q *Queries) ClanFetch(ctx context.Context, cID int64) (Clan, error) {
	row := q.db.QueryRowContext(ctx, clanFetch, cID)
	var i Clan
	err := row.Scan(
		&i.CID,
		&i.CName,
		&i.CPassword,
		&i.CMembers,
		&i.EIDLeader,
		&i.EIDAsst,
		&i.EIDFa1,
		&i.EIDFa2,
		&i.CTitle,
		&i.CUrl,
		&i.CPic,
	)
	return i, err
}

const userAttributesUpdate = `-- name: UserAttributesUpdate :one
UPDATE users
SET u_flags      = ?,
    u_name       = ?,
    u_comment    = ?,
    u_timezone   = ?,
    u_style      = ?,
    u_lang       = ?,
    u_dateformat = ?,
    u_lastip     = ?,
    u_kills      = ?,
    u_deaths     = ?,
    u_offsucc    = ?,
    u_offtotal   = ?,
    u_defsucc    = ?,
    u_deftotal   = ?,
    u_numplays   = ?,
    u_sucplays   = ?,
    u_avgrank    = ?,
    u_bestrank   = ?,
    u_lastdate   = datetime('now')
WHERE u_id = ?
RETURNING u_lastdate
`

type UserAttributesUpdateParams struct {
	UFlags      sql.NullInt64
	UName       sql.NullString
	UComment    sql.NullString
	UTimezone   sql.NullInt64
	UStyle      sql.NullString
	ULang       sql.NullString
	UDateformat sql.NullString
	ULastip     sql.NullString
	UKills      sql.NullInt64
	UDeaths     sql.NullInt64
	UOffsucc    sql.NullInt64
	UOfftotal   sql.NullInt64
	UDefsucc    sql.NullInt64
	UDeftotal   sql.NullInt64
	UNumplays   sql.NullInt64
	USucplays   sql.NullInt64
	UAvgrank    sql.NullFloat64
	UBestrank   sql.NullFloat64
	UID         int64
}

func (q *Queries) UserAttributesUpdate(ctx context.Context, arg UserAttributesUpdateParams) (sql.NullTime, error) {
	row := q.db.QueryRowContext(ctx, userAttributesUpdate,
		arg.UFlags,
		arg.UName,
		arg.UComment,
		arg.UTimezone,
		arg.UStyle,
		arg.ULang,
		arg.UDateformat,
		arg.ULastip,
		arg.UKills,
		arg.UDeaths,
		arg.UOffsucc,
		arg.UOfftotal,
		arg.UDefsucc,
		arg.UDeftotal,
		arg.UNumplays,
		arg.USucplays,
		arg.UAvgrank,
		arg.UBestrank,
		arg.UID,
	)
	var u_lastdate sql.NullTime
	err := row.Scan(&u_lastdate)
	return u_lastdate, err
}

const userCreate = `-- name: UserCreate :one
INSERT INTO users (u_username, u_email, u_createdate, u_lastdate)
VALUES (?, ?, datetime('now'), datetime('now'))
RETURNING u_id, u_createdate, u_lastdate
`

type UserCreateParams struct {
	UUsername string
	UEmail    string
}

type UserCreateRow struct {
	UID         int64
	UCreatedate sql.NullTime
	ULastdate   sql.NullTime
}

func (q *Queries) UserCreate(ctx context.Context, arg UserCreateParams) (UserCreateRow, error) {
	row := q.db.QueryRowContext(ctx, userCreate, arg.UUsername, arg.UEmail)
	var i UserCreateRow
	err := row.Scan(&i.UID, &i.UCreatedate, &i.ULastdate)
	return i, err
}

const userPasswordUpdate = `-- name: UserPasswordUpdate :one
UPDATE users
SET u_password = ?,
    u_lastdate = datetime('now')
WHERE u_id = ?
RETURNING u_lastdate
`

type UserPasswordUpdateParams struct {
	UPassword sql.NullString
	UID       int64
}

func (q *Queries) UserPasswordUpdate(ctx context.Context, arg UserPasswordUpdateParams) (sql.NullTime, error) {
	row := q.db.QueryRowContext(ctx, userPasswordUpdate, arg.UPassword, arg.UID)
	var u_lastdate sql.NullTime
	err := row.Scan(&u_lastdate)
	return u_lastdate, err
}