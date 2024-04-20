// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package orm

import (
	"database/sql"
	"fmt"
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm/sqlc"
	"log"
	"strings"
)

const (
	// Clan relation flags
	CRFLAG_MUTUAL = 0x01 // Clan relation is mutual - set to complete an alliance, clear to stop a war
	CRFLAG_ALLY   = 0x02 // Clan relation describes an alliance
	CRFLAG_WAR    = 0x04 // Clan relation describes a war

	// Clan forum thread flags
	CTFLAG_NEWS   = 0x01 // Topic contains News postings for the clan, visible on main page
	CTFLAG_STICKY = 0x02 // Topic is sticky and appears in bold at the top of the list
	CTFLAG_LOCK   = 0x04 // Topic has been locked - normal members may not post
	CTFLAG_DELETE = 0x08 // Topic has been deleted

	// Clan forum message flags
	CMFLAG_EDIT   = 0x01 // Post has been edited
	CMFLAG_DELETE = 0x02 // Post has been deleted

	// Clan invite flags
	CIFLAG_PERM = 0x01 // Clan invitation is permanent, effectively a whitelist entry

	// Empire flags
	EFLAG_MOD     = 0x0000 // Unused
	EFLAG_ADMIN   = 0x0002 // Empire is owned by moderator/administrator and cannot interact with other empires
	EFLAG_DISABLE = 0x0004 // Empire is disabled
	EFLAG_VALID   = 0x0008 // Empire has submitted their validation code
	EFLAG_DELETE  = 0x0010 // Empire is flagged for deletion
	EFLAG_MULTI   = 0x0020 // Empire is one of multiple accounts being accessed from the same location (legally or not)
	EFLAG_NOTIFY  = 0x0040 // Empire is in a notification state and cannot perform actions (and will not update idle time)
	EFLAG_ONLINE  = 0x0080 // Empire is currently logged in
	EFLAG_SILENT  = 0x0100 // Empire is prohibited from sending private messages to non-Administrators
	EFLAG_LOGGED  = 0x0200 // All actions performed by empire are logged with a special event code

	// Empire message flags
	MFLAG_DELETE = 0x01 // Message has been deleted
	MFLAG_READ   = 0x02 // Message has been read
	MFLAG_REPLY  = 0x04 // Message has been replied to
	MFLAG_REPORT = 0x08 // Message has been reported for abuse
	MFLAG_DEAD   = 0x10 // Message sender is dead

	// Empire news flags
	NFLAG_READ   = 0x01 // News item has been read
	NFLAG_LOCK   = 0x02 // News item is currently being processed
	NFLAG_GOTTEN = 0x04 // Items attached to the news message have been received

	// History round flags
	HRFLAG_CLANS = 0x01 // Round had clans enabled
	HRFLAG_SCORE = 0x02 // Round ranked empires by score rather than networth

	// History empire flags
	HEFLAG_PROTECT = 0x01        // Empire was protected, whether vacation or newly registered
	HEFLAG_ADMIN   = EFLAG_ADMIN // Empire was owned by a moderator/administrator

	// User flags
	UFLAG_MOD     = 0x01 // User has Moderator privileges (can set/clear multi and disabled flags, can browse empire messages)
	UFLAG_ADMIN   = 0x02 // User has Administrator privileges (can grant/revoke privileges, delete/rename empires, login as anyone, edit clans)
	UFLAG_DISABLE = 0x04 // User account is disabled, cannot create new empires (but can still login to existing ones)
	UFLAG_VALID   = 0x08 // User account's email address has been validated at least once
	UFLAG_CLOSED  = 0x10 // User account has been voluntarily closed, cannot create new empires or login to existing ones
	UFLAG_WATCH   = 0x20 // User account is suspected of abuse
)

func (db *DB) UserCreate(userName, email string) (*model.User_t, error) {
	if userName == "" {
		return nil, fmt.Errorf("username must not be blank")
	} else if len(userName) < 6 {
		return nil, fmt.Errorf("username must be at least 6 characters")
	} else if len(userName) >= 255 {
		return nil, fmt.Errorf("username must be less than 255 characters")
	} else if strings.TrimSpace(userName) != userName {
		return nil, fmt.Errorf("username must not start or end with spaces")
	}
	if email == "" {
		return nil, fmt.Errorf("email must not be blank")
	} else if len(email) < 6 {
		return nil, fmt.Errorf("email must be at least 6 characters")
	} else if len(email) >= 255 {
		return nil, fmt.Errorf("email must be less than 255 characters")
	} else if strings.TrimSpace(email) != email {
		return nil, fmt.Errorf("email must not start or end with spaces")
	} else if !isValidEmailAddress(email) {
		return nil, fmt.Errorf("email must parse")
	}

	var user model.User_t
	if row, err := db.db.UserCreate(db.ctx, sqlc.UserCreateParams{
		UUsername: userName,
		UEmail:    email,
	}); err != nil {
		return nil, err
	} else {
		user.Id = int(row.UID)
		user.UserName = userName
		user.Email = email
		user.CreateDate = row.UCreatedate.Time
		user.LastDate = row.ULastdate.Time
	}

	return &user, nil
}

func (db *DB) UserAttributesUpdate(user *model.User_t) error {
	parms := sqlc.UserAttributesUpdateParams{
		UFlags:      sql.NullInt64{Valid: true, Int64: 0},
		UName:       sql.NullString{Valid: true, String: user.Nickname},
		UComment:    sql.NullString{},
		UTimezone:   sql.NullInt64{},
		UStyle:      sql.NullString{},
		ULang:       sql.NullString{},
		UDateformat: sql.NullString{},
		ULastip:     sql.NullString{Valid: true, String: "localhost"},
		UKills:      sql.NullInt64{},
		UDeaths:     sql.NullInt64{},
		UOffsucc:    sql.NullInt64{},
		UOfftotal:   sql.NullInt64{},
		UDefsucc:    sql.NullInt64{},
		UDeftotal:   sql.NullInt64{},
		UNumplays:   sql.NullInt64{},
		USucplays:   sql.NullInt64{},
		UAvgrank:    sql.NullFloat64{},
		UBestrank:   sql.NullFloat64{},
		UID:         int64(user.Id),
	}

	if user.Flags.Admin {
		parms.UFlags.Int64 |= UFLAG_ADMIN
	}
	if user.Flags.Closed {
		parms.UFlags.Int64 |= UFLAG_CLOSED
	}
	if user.Flags.Disabled {
		parms.UFlags.Int64 |= UFLAG_DISABLE
	}
	if user.Flags.Mod {
		parms.UFlags.Int64 |= UFLAG_MOD
	}
	if user.Flags.Valid {
		parms.UFlags.Int64 |= UFLAG_VALID
	}
	if user.Flags.Watch {
		parms.UFlags.Int64 |= UFLAG_WATCH
	}

	if lastDate, err := db.db.UserAttributesUpdate(db.ctx, parms); err != nil {
		return err
	} else {
		user.LastDate = lastDate.Time
	}

	return nil
}

func (db *DB) UserPasswordUpdate(user *model.User_t) error {
	log.Printf("orm: userPasswordUpdate: please bcrypt the password!\n")
	parms := sqlc.UserPasswordUpdateParams{
		UPassword: sql.NullString{Valid: true, String: user.Password},
		UID:       int64(user.Id),
	}
	if lastDate, err := db.db.UserPasswordUpdate(db.ctx, parms); err != nil {
		return err
	} else {
		user.LastDate = lastDate.Time
	}
	return nil
}

func isValidEmailAddress(address string) bool {
	if len(address) < 6 || len(address) > 255 {
		return false
	} else if strings.Count(address, "@") != 1 {
		return false
	} else if pos := strings.Index(address, "@"); pos < 3 || pos > len(address)-3 {
		return false
	}
	return true
}
