// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package orm

import (
	"database/sql"
	"fmt"
	"github.com/mdhender/promisance/app/model"
	"github.com/mdhender/promisance/app/orm/sqlc"
	"log"
	"strings"
	"time"
)

const (
	// Clan relation flags
	CRFLAG_ALLY   = 0x02 // Clan relation describes an alliance
	CRFLAG_MUTUAL = 0x01 // Clan relation is mutual - set to complete an alliance, clear to stop a war
	CRFLAG_WAR    = 0x04 // Clan relation describes a war

	// Clan forum thread flags
	CTFLAG_DELETE = 0x08 // Topic has been deleted
	CTFLAG_LOCK   = 0x04 // Topic has been locked - normal members may not post
	CTFLAG_NEWS   = 0x01 // Topic contains News postings for the clan, visible on main page
	CTFLAG_STICKY = 0x02 // Topic is sticky and appears in bold at the top of the list

	// Clan forum message flags
	CMFLAG_DELETE = 0x02 // Post has been deleted
	CMFLAG_EDIT   = 0x01 // Post has been edited

	// Clan invite flags
	CIFLAG_PERM = 0x01 // Clan invitation is permanent, effectively a whitelist entry

	// Empire flags
	EFLAG_ADMIN   = 0x0002 // Empire is owned by moderator/administrator and cannot interact with other empires
	EFLAG_DELETE  = 0x0010 // Empire is flagged for deletion
	EFLAG_DISABLE = 0x0004 // Empire is disabled
	EFLAG_LOGGED  = 0x0200 // All actions performed by empire are logged with a special event code
	EFLAG_MOD     = 0x0000 // Unused
	EFLAG_MULTI   = 0x0020 // Empire is one of multiple accounts being accessed from the same location (legally or not)
	EFLAG_NOTIFY  = 0x0040 // Empire is in a notification state and cannot perform actions (and will not update idle time)
	EFLAG_ONLINE  = 0x0080 // Empire is currently logged in
	EFLAG_SILENT  = 0x0100 // Empire is prohibited from sending private messages to non-Administrators
	EFLAG_VALID   = 0x0008 // Empire has submitted their validation code

	// Empire message flags
	MFLAG_DELETE = 0x01 // Message has been deleted
	MFLAG_READ   = 0x02 // Message has been read
	MFLAG_REPLY  = 0x04 // Message has been replied to
	MFLAG_REPORT = 0x08 // Message has been reported for abuse
	MFLAG_DEAD   = 0x10 // Message sender is dead

	// Empire news flags
	NFLAG_GOTTEN = 0x04 // Items attached to the news message have been received
	NFLAG_LOCK   = 0x02 // News item is currently being processed
	NFLAG_READ   = 0x01 // News item has been read

	// History round flags
	HRFLAG_CLANS = 0x01 // Round had clans enabled
	HRFLAG_SCORE = 0x02 // Round ranked empires by score rather than networth

	// History empire flags
	HEFLAG_ADMIN   = EFLAG_ADMIN // Empire was owned by a moderator/administrator
	HEFLAG_PROTECT = 0x01        // Empire was protected, whether vacation or newly registered

	RACE_DROW    = 8
	RACE_DWARF   = 3
	RACE_ELF     = 2
	RACE_GNOME   = 5
	RACE_GOBLIN  = 9
	RACE_GREMLIN = 6
	RACE_HUMAN   = 1
	RACE_ORC     = 7
	RACE_TROLL   = 4

	// User flags
	UFLAG_ADMIN   = 0x02 // User has Administrator privileges (can grant/revoke privileges, delete/rename empires, login as anyone, edit clans)
	UFLAG_CLOSED  = 0x10 // User account has been voluntarily closed, cannot create new empires or login to existing ones
	UFLAG_DISABLE = 0x04 // User account is disabled, cannot create new empires (but can still login to existing ones)
	UFLAG_MOD     = 0x01 // User has Moderator privileges (can set/clear multi and disabled flags, can browse empire messages)
	UFLAG_VALID   = 0x08 // User account's email address has been validated at least once
	UFLAG_WATCH   = 0x20 // User account is suspected of abuse
)

func (db *DB) AuthenticatedUserFetch(username, password string) (*model.User_t, error) {
	if username == "" || password == "" {
		return nil, sql.ErrNoRows
	}
	row, err := db.db.AuthenticatedUserFetch(db.ctx, sqlc.AuthenticatedUserFetchParams{
		UUsername: username,
		UPassword: sql.NullString{Valid: true, String: password},
	})
	if err != nil {
		return nil, err
	}
	user := &model.User_t{
		Id:       int(row.UID),
		UserName: row.UUsername,
		Flags:    intToUserFlags(row.UFlags),
	}
	if row.UComment.Valid {
		user.Comment = row.UComment.String
	}
	return user, nil
}

func (db *DB) EmpireActiveCount() (int, error) {
	count, err := db.db.EmpireActiveUserCount(db.ctx)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (db *DB) EmpireCreate(user *model.User_t, name string, race string) (*model.Empire_t, error) {
	var raceFlag int64
	switch race {
	case "HUMAN":
		raceFlag = RACE_HUMAN
	default:
		return nil, fmt.Errorf("unknown race: %s", race)
	}

	var empire model.Empire_t
	if id, err := db.db.EmpireCreate(db.ctx, sqlc.EmpireCreateParams{
		UID:   int64(user.Id),
		EName: name,
		ERace: raceFlag,
	}); err != nil {
		return nil, err
	} else {
		empire.Id = int(id)
		empire.UserId = user.Id
		empire.Name = name
		empire.Race = int(raceFlag)
	}

	return &empire, nil
}

func (db *DB) EmpireAttributesUpdate(empire *model.Empire_t) error {
	parms := sqlc.EmpireAttributesUpdateParams{
		EID:          int64(empire.Id),
		EFlags:       empireFlagsToInt(empire.Flags),
		EValcode:     sql.NullString{Valid: true, String: empire.ValCode},
		EReason:      sql.NullString{Valid: true, String: empire.Reason},
		EVacation:    sql.NullInt64{Valid: true, Int64: int64(empire.Vacation)},
		EIdle:        sql.NullInt64{Valid: true, Int64: int64(empire.Idle)},
		EEra:         sql.NullInt64{Valid: true, Int64: int64(empire.Era)},
		ERank:        sql.NullInt64{Valid: true, Int64: int64(empire.Rank)},
		ESharing:     sql.NullInt64{Valid: true, Int64: int64(empire.Sharing)},
		EAttacks:     sql.NullInt64{Valid: true, Int64: int64(empire.Attacks)},
		EOffsucc:     sql.NullInt64{Valid: true, Int64: int64(empire.OffSucc)},
		EOfftotal:    sql.NullInt64{Valid: true, Int64: int64(empire.OffTotal)},
		EDefsucc:     sql.NullInt64{Valid: true, Int64: int64(empire.DefSucc)},
		EDeftotal:    sql.NullInt64{Valid: true, Int64: int64(empire.DefTotal)},
		EKills:       sql.NullInt64{Valid: true, Int64: int64(empire.Kills)},
		EScore:       sql.NullInt64{Valid: true, Int64: int64(empire.Score)},
		EKilledby:    sql.NullInt64{Valid: true, Int64: int64(empire.KilledBy)},
		EKillclan:    sql.NullInt64{Valid: true, Int64: int64(empire.KillClan)},
		ETurns:       sql.NullInt64{Valid: true, Int64: int64(empire.Turns)},
		EStoredturns: sql.NullInt64{Valid: true, Int64: int64(empire.StoredTurns)},
		ETurnsused:   sql.NullInt64{Valid: true, Int64: int64(empire.TurnsUsed)},
		ENetworth:    sql.NullInt64{Valid: true, Int64: int64(empire.NetWorth)},
		ECash:        sql.NullInt64{Valid: true, Int64: int64(empire.Cash)},
		EFood:        sql.NullInt64{Valid: true, Int64: int64(empire.Food)},
		EPeasants:    sql.NullInt64{Valid: true, Int64: int64(empire.Peasants)},
		ETrparm:      sql.NullInt64{Valid: true, Int64: int64(empire.TrpArm)},
		ETrplnd:      sql.NullInt64{Valid: true, Int64: int64(empire.TrpLnd)},
		ETrpfly:      sql.NullInt64{Valid: true, Int64: int64(empire.TrpFly)},
		ETrpsea:      sql.NullInt64{Valid: true, Int64: int64(empire.TrpSea)},
		ETrpwiz:      sql.NullInt64{Valid: true, Int64: int64(empire.TrpWiz)},
		EHealth:      sql.NullInt64{Valid: true, Int64: int64(empire.Health)},
		ERunes:       sql.NullInt64{Valid: true, Int64: int64(empire.Runes)},
		EIndarm:      sql.NullInt64{Valid: true, Int64: int64(empire.IndArm)},
		EIndlnd:      sql.NullInt64{Valid: true, Int64: int64(empire.IndLnd)},
		EIndfly:      sql.NullInt64{Valid: true, Int64: int64(empire.IndFly)},
		EIndsea:      sql.NullInt64{Valid: true, Int64: int64(empire.IndSea)},
		ELand:        sql.NullInt64{Valid: true, Int64: int64(empire.Land)},
		EBldpop:      sql.NullInt64{Valid: true, Int64: int64(empire.BldPop)},
		EBldcash:     sql.NullInt64{Valid: true, Int64: int64(empire.BldCash)},
		EBldtrp:      sql.NullInt64{Valid: true, Int64: int64(empire.BldTrp)},
		EBldcost:     sql.NullInt64{Valid: true, Int64: int64(empire.BldCost)},
		EBldwiz:      sql.NullInt64{Valid: true, Int64: int64(empire.BldWiz)},
		EBldfood:     sql.NullInt64{Valid: true, Int64: int64(empire.BldFood)},
		EBlddef:      sql.NullInt64{Valid: true, Int64: int64(empire.BldDef)},
		EFreeland:    sql.NullInt64{Valid: true, Int64: int64(empire.Freeland)},
		ETax:         sql.NullInt64{Valid: true, Int64: int64(empire.Tax)},
		EBank:        sql.NullInt64{Valid: true, Int64: int64(empire.Bank)},
		ELoan:        sql.NullInt64{Valid: true, Int64: int64(empire.Loan)},
		EMktarm:      sql.NullInt64{Valid: true, Int64: int64(empire.MktArm)},
		EMktlnd:      sql.NullInt64{Valid: true, Int64: int64(empire.MktLnd)},
		EMktfly:      sql.NullInt64{Valid: true, Int64: int64(empire.MktFly)},
		EMktsea:      sql.NullInt64{Valid: true, Int64: int64(empire.MktSea)},
		EMktfood:     sql.NullInt64{Valid: true, Int64: int64(empire.MktFood)},
		EMktperarm:   sql.NullInt64{Valid: true, Int64: int64(empire.MktPerArm)},
		EMktperlnd:   sql.NullInt64{Valid: true, Int64: int64(empire.MktPerLnd)},
		EMktperfly:   sql.NullInt64{Valid: true, Int64: int64(empire.MktPerFly)},
		EMktpersea:   sql.NullInt64{Valid: true, Int64: int64(empire.MktPerSea)},
	}
	return db.db.EmpireAttributesUpdate(db.ctx, parms)
}

func (db *DB) EmpireFetch(id int) (*model.Empire_t, error) {
	row, err := db.db.EmpireFetch(db.ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return &model.Empire_t{
		Id:          int(row.EID),
		UserId:      int(row.UID),
		OldUserId:   nvlInt(row.UOldid),
		SignupDate:  nvlTime(row.ESignupdate),
		Flags:       intToEmpireFlags(row.EFlags),
		ValCode:     nvlString(row.EValcode),
		Reason:      nvlString(row.EReason),
		Vacation:    nvlInt(row.EVacation),
		Idle:        nvlInt(row.EIdle),
		Name:        row.EName,
		Race:        int(row.ERace),
		Era:         nvlInt(row.EEra),
		Rank:        nvlInt(row.ERank),
		CId:         nvlInt(row.CID),
		OldCId:      nvlInt(row.COldid),
		Sharing:     nvlInt(row.ESharing),
		Attacks:     nvlInt(row.EAttacks),
		OffSucc:     nvlInt(row.EOffsucc),
		OffTotal:    nvlInt(row.EOfftotal),
		DefSucc:     nvlInt(row.EDefsucc),
		DefTotal:    nvlInt(row.EDeftotal),
		Kills:       nvlInt(row.EKills),
		Score:       nvlInt(row.EScore),
		KilledBy:    nvlInt(row.EKilledby),
		KillClan:    nvlInt(row.EKillclan),
		Turns:       nvlInt(row.ETurns),
		StoredTurns: nvlInt(row.EStoredturns),
		TurnsUsed:   nvlInt(row.ETurnsused),
		NetWorth:    nvlInt(row.ENetworth),
		Cash:        nvlInt(row.ECash),
		Food:        nvlInt(row.EFood),
		Peasants:    nvlInt(row.EPeasants),
		TrpArm:      nvlInt(row.ETrparm),
		TrpLnd:      nvlInt(row.ETrplnd),
		TrpFly:      nvlInt(row.ETrpfly),
		TrpSea:      nvlInt(row.ETrpsea),
		TrpWiz:      nvlInt(row.ETrpwiz),
		Health:      nvlInt(row.EHealth),
		Runes:       nvlInt(row.ERunes),
		IndArm:      nvlInt(row.EIndarm),
		IndLnd:      nvlInt(row.EIndlnd),
		IndFly:      nvlInt(row.EIndfly),
		IndSea:      nvlInt(row.EIndsea),
		Land:        nvlInt(row.ELand),
		BldPop:      nvlInt(row.EBldpop),
		BldCash:     nvlInt(row.EBldcash),
		BldTrp:      nvlInt(row.EBldtrp),
		BldCost:     nvlInt(row.EBldcost),
		BldWiz:      nvlInt(row.EBldwiz),
		BldFood:     nvlInt(row.EBldfood),
		BldDef:      nvlInt(row.EBlddef),
		Freeland:    nvlInt(row.EFreeland),
		Tax:         nvlInt(row.ETax),
		Bank:        nvlInt(row.EBank),
		Loan:        nvlInt(row.ELoan),
		MktArm:      nvlInt(row.EMktarm),
		MktLnd:      nvlInt(row.EMktlnd),
		MktFly:      nvlInt(row.EMktfly),
		MktSea:      nvlInt(row.EMktsea),
		MktFood:     nvlInt(row.EMktfood),
		MktPerArm:   nvlInt(row.EMktperarm),
		MktPerLnd:   nvlInt(row.EMktperlnd),
		MktPerFly:   nvlInt(row.EMktperfly),
		MktPerSea:   nvlInt(row.EMktpersea),
	}, nil
}

func (db *DB) EmpireUpdateFlags(empire *model.Empire_t) error {
	return db.db.EmpireUpdateFlags(db.ctx, sqlc.EmpireUpdateFlagsParams{
		EFlags: empireFlagsToInt(empire.Flags),
		EID:    int64(empire.Id),
	})
}

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

func (db *DB) UserAccessUpdate(user *model.User_t) error {
	parms := sqlc.UserAccessUpdateParams{
		ULastip: sql.NullString{Valid: true, String: user.LastIP},
		UID:     int64(user.Id),
	}
	if lastDate, err := db.db.UserAccessUpdate(db.ctx, parms); err != nil {
		return err
	} else {
		user.LastDate = lastDate.Time
	}
	return nil
}

func (db *DB) UserActiveEmpires(userId int) ([]*model.Empire_t, error) {
	rows, err := db.db.UserActiveEmpires(db.ctx, sqlc.UserActiveEmpiresParams{
		UID:    int64(userId),
		EFlags: sql.NullInt64{Valid: true, Int64: EFLAG_DELETE},
	})
	if err != nil {
		return nil, err
	}
	var empires []*model.Empire_t
	for _, row := range rows {
		empires = append(empires, &model.Empire_t{Id: int(row.EID), Flags: intToEmpireFlags(row.EFlags)})
	}
	return empires, nil
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

func (db *DB) UserFetch(id int) (*model.User_t, error) {
	row, err := db.db.UserFetch(db.ctx, int64(id))
	if err != nil {
		return nil, err
	}
	user := &model.User_t{
		Id:         int(row.UID),
		UserName:   row.UUsername,
		Flags:      intToUserFlags(row.UFlags),
		Nickname:   nvlString(row.UName),
		Email:      row.UEmail,
		Comment:    nvlString(row.UComment),
		TimeZone:   nvlInt(row.UTimezone),
		Style:      nvlString(row.UStyle),
		Lang:       nvlString(row.ULang),
		DateFormat: nvlString(row.UDateformat),
		LastIP:     nvlString(row.ULastip),
		Kills:      nvlInt(row.UKills),
		Deaths:     nvlInt(row.UDeaths),
		OffSucc:    nvlInt(row.UOffsucc),
		OffTotal:   nvlInt(row.UOfftotal),
		DefSucc:    nvlInt(row.UDefsucc),
		DefTotal:   nvlInt(row.UDeftotal),
		NumPlays:   nvlInt(row.UNumplays),
		SucPlays:   nvlInt(row.USucplays),
		AvgRank:    nvlFloat(row.UAvgrank),
		Bestrank:   nvlFloat(row.UBestrank),
		CreateDate: nvlTime(row.UCreatedate),
		LastDate:   nvlTime(row.ULastdate),
	}
	return user, nil
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

func (db *DB) WorldVarsFetch() (*model.World_t, error) {
	row, err := db.db.WorldVarsFetch(db.ctx)
	if err != nil {
		return nil, err
	}
	return &model.World_t{
		Id:                    int(row.WvID),
		LottoCurrentJackpot:   int(row.LottoCurrentJackpot),
		LottoYesterdayJackpot: int(row.LottoYesterdayJackpot),
		LottoLastPicked:       int(row.LottoLastPicked),
		LottoLastWinner:       int(row.LottoLastWinner),
		LottoJackpotIncrease:  int(row.LottoJackpotIncrease),
		RoundTimeBegin:        row.RoundTimeBegin,
		RoundTimeClosing:      row.RoundTimeClosing,
		RoundTimeEnd:          row.RoundTimeEnd,
		TurnsNext:             row.TurnsNext,
		TurnsNextHourly:       row.TurnsNextHourly,
		TurnsNextDaily:        row.TurnsNextDaily,
	}, nil
}

func (db *DB) WorldVarsInitialize(world *model.World_t) error {
	parms := sqlc.WorldVarsInitializeParams{
		LottoCurrentJackpot:   int64(world.LottoCurrentJackpot),
		LottoYesterdayJackpot: int64(world.LottoYesterdayJackpot),
		LottoLastPicked:       int64(world.LottoLastPicked),
		LottoLastWinner:       int64(world.LottoLastWinner),
		LottoJackpotIncrease:  int64(world.LottoJackpotIncrease),
		RoundTimeBegin:        world.RoundTimeBegin,
		RoundTimeClosing:      world.RoundTimeClosing,
		RoundTimeEnd:          world.RoundTimeEnd,
		TurnsNext:             world.TurnsNext,
		TurnsNextHourly:       world.TurnsNextHourly,
		TurnsNextDaily:        world.TurnsNextDaily,
	}
	return db.db.WorldVarsInitialize(db.ctx, parms)
}

func (db *DB) WorldVarsUpdate(world *model.World_t) error {
	parms := sqlc.WorldVarsUpdateParams{
		LottoCurrentJackpot:   int64(world.LottoCurrentJackpot),
		LottoYesterdayJackpot: int64(world.LottoYesterdayJackpot),
		LottoLastPicked:       int64(world.LottoLastPicked),
		LottoLastWinner:       int64(world.LottoLastWinner),
		LottoJackpotIncrease:  int64(world.LottoJackpotIncrease),
		RoundTimeBegin:        world.RoundTimeBegin,
		RoundTimeClosing:      world.RoundTimeClosing,
		RoundTimeEnd:          world.RoundTimeEnd,
		TurnsNext:             world.TurnsNext,
		TurnsNextHourly:       world.TurnsNextHourly,
		TurnsNextDaily:        world.TurnsNextDaily,
	}
	return db.db.WorldVarsUpdate(db.ctx, parms)
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

func empireFlagsToInt(flags model.EmpireFlag_t) sql.NullInt64 {
	var bits int
	if flags.Admin {
		bits |= EFLAG_ADMIN
	}
	if flags.Delete {
		bits |= EFLAG_DELETE
	}
	if flags.Disable {
		bits |= EFLAG_DISABLE
	}
	if flags.Logged {
		bits |= EFLAG_LOGGED
	}
	if flags.Mod {
		bits |= EFLAG_MOD
	}
	if flags.Multi {
		bits |= EFLAG_MULTI
	}
	if flags.Notify {
		bits |= EFLAG_NOTIFY
	}
	if flags.Online {
		bits |= EFLAG_ONLINE
	}
	if flags.Silent {
		bits |= EFLAG_SILENT
	}
	if flags.Valid {
		bits |= EFLAG_VALID
	}
	return sql.NullInt64{Valid: true, Int64: int64(bits)}
}

func intToEmpireFlags(flags sql.NullInt64) model.EmpireFlag_t {
	var bits int
	if flags.Valid {
		bits = int(flags.Int64)
	}
	return model.EmpireFlag_t{
		Admin:   (bits & EFLAG_ADMIN) != 0,
		Delete:  (bits & EFLAG_DELETE) != 0,
		Disable: (bits & EFLAG_DISABLE) != 0,
		Logged:  (bits & EFLAG_LOGGED) != 0,
		Mod:     (bits & EFLAG_MOD) != 0,
		Multi:   (bits & EFLAG_MULTI) != 0,
		Notify:  (bits & EFLAG_NOTIFY) != 0,
		Online:  (bits & EFLAG_ONLINE) != 0,
		Silent:  (bits & EFLAG_SILENT) != 0,
		Valid:   (bits & EFLAG_VALID) != 0,
	}
}

func userFlagsToInt(flags model.UserFlag_t) sql.NullInt64 {
	var bits int
	if flags.Admin {
		bits |= UFLAG_ADMIN
	}
	if flags.Closed {
		bits |= UFLAG_CLOSED
	}
	if flags.Disabled {
		bits |= UFLAG_DISABLE
	}
	if flags.Mod {
		bits |= UFLAG_MOD
	}
	if flags.Valid {
		bits |= UFLAG_VALID
	}
	if flags.Watch {
		bits |= UFLAG_WATCH
	}
	return sql.NullInt64{Valid: true, Int64: int64(bits)}
}

func intToUserFlags(flags sql.NullInt64) model.UserFlag_t {
	var bits int
	if flags.Valid {
		bits = int(flags.Int64)
	}
	return model.UserFlag_t{
		Admin:    (bits & UFLAG_ADMIN) != 0,
		Closed:   (bits & UFLAG_CLOSED) != 0,
		Disabled: (bits & UFLAG_DISABLE) != 0,
		Mod:      (bits & UFLAG_MOD) != 0,
		Valid:    (bits & UFLAG_VALID) != 0,
		Watch:    (bits & UFLAG_WATCH) != 0,
	}
}

func nvlFloat(v sql.NullFloat64) float64 {
	if !v.Valid {
		return 0
	}
	return v.Float64
}

func nvlInt(v sql.NullInt64) int {
	if !v.Valid {
		return 0
	}
	return int(v.Int64)
}

func nvlString(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

func nvlTime(v sql.NullTime) time.Time {
	if !v.Valid {
		return time.Time{}
	}
	return v.Time
}
