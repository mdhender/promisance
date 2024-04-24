// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package model

import (
	"time"
)

type Empire_t struct {
	Id          int
	UserId      int
	OldUserId   int
	SignupDate  time.Time
	Flags       EmpireFlag_t
	ValCode     string
	Reason      string
	Vacation    int
	Idle        int
	Name        string
	Race        int
	Era         int
	Rank        int
	CId         int
	OldCId      int
	Sharing     int
	Attacks     int
	OffSucc     int
	OffTotal    int
	DefSucc     int
	DefTotal    int
	Kills       int
	Score       int
	KilledBy    int
	KillClan    int
	Turns       int
	StoredTurns int
	TurnsUsed   int
	NetWorth    int
	Cash        int
	Food        int
	Peasants    int
	TrpArm      int
	TrpLnd      int
	TrpFly      int
	TrpSea      int
	TrpWiz      int
	Health      int
	Runes       int
	IndArm      int
	IndLnd      int
	IndFly      int
	IndSea      int
	Land        int
	BldPop      int
	BldCash     int
	BldTrp      int
	BldCost     int
	BldWiz      int
	BldFood     int
	BldDef      int
	Freeland    int
	Tax         int
	Bank        int
	Loan        int
	MktArm      int
	MktLnd      int
	MktFly      int
	MktSea      int
	MktFood     int
	MktPerArm   int
	MktPerLnd   int
	MktPerFly   int
	MktPerSea   int
}

type EmpireFlag_t struct {
	// Unused
	Mod bool
	// Empire is owned by moderator/administrator and cannot interact with other empires
	Admin bool
	// Empire is disabled
	Disable bool
	// Empire has submitted their validation code
	Valid bool
	// Empire is flagged for deletion
	Delete bool
	// Empire is one of multiple accounts being accessed from the same location (legally or not)
	Multi bool
	// Empire is in a notification state and cannot perform actions (and will not update idle time)
	Notify bool
	// Empire is currently logged in
	Online bool
	// Empire is prohibited from sending private messages to non-Administrators
	Silent bool
	// All actions performed by empire are logged with a special event code
	Logged bool
}

type RoundData_t struct {
	Signup     bool
	Started    bool
	Closing    bool
	Finished   bool
	TimeNotice string
}

type User_t struct {
	Id         int
	UserName   string
	Password   string
	Flags      UserFlag_t
	Nickname   string
	Email      string
	Comment    string
	TimeZone   int
	Style      string
	Lang       string
	DateFormat string
	LastIP     string
	Kills      int
	Deaths     int
	OffSucc    int
	OffTotal   int
	DefSucc    int
	DefTotal   int
	NumPlays   int
	SucPlays   int
	AvgRank    float64
	Bestrank   float64
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

type World_t struct {
	Id                    int
	LottoCurrentJackpot   int
	LottoYesterdayJackpot int
	LottoLastPicked       int
	LottoLastWinner       int
	LottoJackpotIncrease  int
	RoundTimeBegin        time.Time
	RoundTimeClosing      time.Time
	RoundTimeEnd          time.Time
	TurnsNext             time.Time
	TurnsNextHourly       time.Time
	TurnsNextDaily        time.Time
}
