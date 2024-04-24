// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"database/sql"
	"time"
)

type Clan struct {
	CID       int64
	CName     string
	CPassword string
	CMembers  int64
	EIDLeader int64
	EIDAsst   int64
	EIDFa1    int64
	EIDFa2    int64
	CTitle    string
	CUrl      string
	CPic      string
}

type ClanInvite struct {
	CiID    int64
	CID     interface{}
	EID1    interface{}
	EID2    interface{}
	CiFlags interface{}
	CiTime  int64
}

type ClanMessage struct {
	CmID    int64
	CtID    interface{}
	EID     interface{}
	CmBody  string
	CmTime  int64
	CmFlags interface{}
}

type ClanNews struct {
	CnID    int64
	CnTime  int64
	CID     interface{}
	EID1    interface{}
	CID2    interface{}
	EID2    interface{}
	CnEvent interface{}
}

type ClanRelation struct {
	CrID    int64
	CID1    interface{}
	CID2    interface{}
	CrFlags interface{}
	CrTime  int64
}

type ClanTopic struct {
	CtID      int64
	CID       interface{}
	CtSubject string
	CtFlags   interface{}
}

type Empire struct {
	EID          int64
	UID          int64
	UOldid       sql.NullInt64
	ESignupdate  sql.NullTime
	EFlags       sql.NullInt64
	EValcode     sql.NullString
	EReason      sql.NullString
	EVacation    sql.NullInt64
	EIdle        sql.NullInt64
	EName        string
	ERace        int64
	EEra         sql.NullInt64
	ERank        sql.NullInt64
	CID          sql.NullInt64
	COldid       sql.NullInt64
	ESharing     sql.NullInt64
	EAttacks     sql.NullInt64
	EOffsucc     sql.NullInt64
	EOfftotal    sql.NullInt64
	EDefsucc     sql.NullInt64
	EDeftotal    sql.NullInt64
	EKills       sql.NullInt64
	EScore       sql.NullInt64
	EKilledby    sql.NullInt64
	EKillclan    sql.NullInt64
	ETurns       sql.NullInt64
	EStoredturns sql.NullInt64
	ETurnsused   sql.NullInt64
	ENetworth    sql.NullInt64
	ECash        sql.NullInt64
	EFood        sql.NullInt64
	EPeasants    sql.NullInt64
	ETrparm      sql.NullInt64
	ETrplnd      sql.NullInt64
	ETrpfly      sql.NullInt64
	ETrpsea      sql.NullInt64
	ETrpwiz      sql.NullInt64
	EHealth      sql.NullInt64
	ERunes       sql.NullInt64
	EIndarm      sql.NullInt64
	EIndlnd      sql.NullInt64
	EIndfly      sql.NullInt64
	EIndsea      sql.NullInt64
	ELand        sql.NullInt64
	EBldpop      sql.NullInt64
	EBldcash     sql.NullInt64
	EBldtrp      sql.NullInt64
	EBldcost     sql.NullInt64
	EBldwiz      sql.NullInt64
	EBldfood     sql.NullInt64
	EBlddef      sql.NullInt64
	EFreeland    sql.NullInt64
	ETax         sql.NullInt64
	EBank        sql.NullInt64
	ELoan        sql.NullInt64
	EMktarm      sql.NullInt64
	EMktlnd      sql.NullInt64
	EMktfly      sql.NullInt64
	EMktsea      sql.NullInt64
	EMktfood     sql.NullInt64
	EMktperarm   sql.NullInt64
	EMktperlnd   sql.NullInt64
	EMktperfly   sql.NullInt64
	EMktpersea   sql.NullInt64
}

type EmpireEffect struct {
	EID     interface{}
	EfName  interface{}
	EfValue int64
}

type EmpireMessage struct {
	MID      int64
	MIDRef   interface{}
	MTime    int64
	EIDSrc   interface{}
	EIDDst   interface{}
	MSubject string
	MBody    string
	MFlags   interface{}
}

type EmpireNews struct {
	NID    int64
	NTime  int64
	EIDSrc interface{}
	CIDSrc interface{}
	EIDDst interface{}
	CIDDst interface{}
	NEvent interface{}
	ND0    int64
	ND1    int64
	ND2    int64
	ND3    int64
	ND4    int64
	ND5    int64
	ND6    int64
	ND7    int64
	ND8    int64
	NFlags interface{}
}

type HistoryClan struct {
	HrID       int64
	HcID       interface{}
	HcMembers  int64
	HcName     string
	HcTitle    string
	HcTotalnet interface{}
}

type HistoryEmpire struct {
	HrID       int64
	HeFlags    interface{}
	UID        interface{}
	HeID       interface{}
	HeName     string
	HeRace     string
	HeEra      string
	HcID       interface{}
	HeOffsucc  interface{}
	HeOfftotal interface{}
	HeDefsucc  interface{}
	HeDeftotal interface{}
	HeKills    interface{}
	HeScore    int64
	HeNetworth interface{}
	HeLand     interface{}
	HeRank     interface{}
}

type HistoryRound struct {
	HrID             int64
	HrName           string
	HrDescription    string
	HrStartdate      int64
	HrStopdate       int64
	HrFlags          interface{}
	HrSmallclansize  interface{}
	HrSmallclans     interface{}
	HrAllclans       interface{}
	HrNonclanempires interface{}
	HrLiveempires    interface{}
	HrDeadempires    interface{}
	HrDelempires     interface{}
	HrAllempires     interface{}
}

type Lock struct {
	LockID interface{}
}

type Log struct {
	LogID     int64
	LogTime   interface{}
	LogType   interface{}
	LogIp     string
	LogPage   string
	LogAction string
	LogLocks  string
	LogText   string
	UID       interface{}
	EID       interface{}
	CID       interface{}
}

type Lottery struct {
	EID     interface{}
	LTicket interface{}
	LCash   interface{}
}

type Market struct {
	KID    int64
	KType  interface{}
	EID    interface{}
	KAmt   interface{}
	KPrice interface{}
	KTime  int64
}

type Permission struct {
	PID         int64
	PType       interface{}
	PCriteria   string
	PComment    string
	PReason     string
	PCreatetime interface{}
	PUpdatetime interface{}
	PLasthit    interface{}
	PHitcount   interface{}
	PExpire     interface{}
}

type Session struct {
	SessID        string
	SessExpiresAt time.Time
	SessUid       int64
	SessEid       int64
}

type Turnlog struct {
	TurnID       int64
	TurnTime     interface{}
	TurnTicks    interface{}
	TurnInterval interface{}
	TurnType     interface{}
	TurnText     string
}

type User struct {
	UID         int64
	UUsername   string
	UPassword   sql.NullString
	UFlags      sql.NullInt64
	UName       sql.NullString
	UEmail      string
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
	UCreatedate sql.NullTime
	ULastdate   sql.NullTime
}

type Var struct {
	VName  string
	VValue string
}

type VarAdjust struct {
	VName   interface{}
	VOffset int64
}

type WorldVar struct {
	WvID                  int64
	LottoCurrentJackpot   int64
	LottoYesterdayJackpot int64
	LottoLastPicked       int64
	LottoLastWinner       int64
	LottoJackpotIncrease  int64
	RoundTimeBegin        time.Time
	RoundTimeClosing      time.Time
	RoundTimeEnd          time.Time
	TurnsNext             time.Time
	TurnsNextHourly       time.Time
	TurnsNextDaily        time.Time
}
