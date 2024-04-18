// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import "time"

const (
	// Entity type IDs, used for locking
	ENT_USER   = 1 // User account
	ENT_EMPIRE = 2 // Empire
	ENT_CLAN   = 3 // Clan
	ENT_VARS   = 4 // World variables
	ENT_MARKET = 5 // Public market items

	// Permission flags
	PERM_EXCEPT = 0x01 // Permission entry is an exception rather than a ban
	PERM_IPV4   = 0x00 // Permission specifies an IPv4 address+mask
	PERM_EMAIL  = 0x02 // Permission specifies an email address mask
	PERM_IPV6   = 0x04 // Permission specifies an IPv6 address+mask
	PERM_MASK   = 0x06 // Bitmask for permission types

	// Lock owner IDs for special functions - used only for potential logging purposes
	LOCK_SCRIPT  = 2147483643 // Utility script
	LOCK_HISTORY = 2147483644 // Record history
	LOCK_RESET   = 2147483645 // Round reset
	LOCK_NEW     = 2147483646 // New entity creation
	LOCK_TURNS   = 2147483647 // Turns script

	// User flags
	UFLAG_MOD     = 0x01 // User has Moderator privileges (can set/clear multi and disabled flags, can browse empire messages)
	UFLAG_ADMIN   = 0x02 // User has Administrator privileges (can grant/revoke privileges, delete/rename empires, login as anyone, edit clans)
	UFLAG_DISABLE = 0x04 // User account is disabled, cannot create new empires (but can still login to existing ones)
	UFLAG_VALID   = 0x08 // User account's email address has been validated at least once
	UFLAG_CLOSED  = 0x10 // User account has been voluntarily closed, cannot create new empires or login to existing ones
	UFLAG_WATCH   = 0x20 // User account is suspected of abuse

	// Empire flags
	//EFLAG_MOD = 0  // Unused
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

	// History round flags
	HRFLAG_CLANS = 0x01 // Round had clans enabled
	HRFLAG_SCORE = 0x02 // Round ranked empires by score rather than networth

	// History empire flags
	HEFLAG_PROTECT = 0x01        // Empire was protected, whether vacation or newly registered
	HEFLAG_ADMIN   = EFLAG_ADMIN // Empire was owned by a moderator/administrator

	// Turn log entry types
	TURN_EVENT = 0 // Normal turn log entry
	TURN_START = 1 // Start of a turn run
	TURN_END   = 2 // End of a turn run
	TURN_ABORT = 3 // Turn run was aborted due to there being nothing to do

	// Public market item types
	MARKET_TRPARM = 0
	MARKET_TRPLND = 1
	MARKET_TRPFLY = 2
	MARKET_TRPSEA = 3
	MARKET_FOOD   = 4

	// Database table names
	CLAN_TABLE           = TABLE_PREFIX + "clan"
	CLAN_INVITE_TABLE    = TABLE_PREFIX + "clan_invite"
	CLAN_MESSAGE_TABLE   = TABLE_PREFIX + "clan_message"
	CLAN_NEWS_TABLE      = TABLE_PREFIX + "clan_news"
	CLAN_RELATION_TABLE  = TABLE_PREFIX + "clan_relation"
	CLAN_TOPIC_TABLE     = TABLE_PREFIX + "clan_topic"
	EMPIRE_TABLE         = TABLE_PREFIX + "empire"
	EMPIRE_EFFECT_TABLE  = TABLE_PREFIX + "empire_effect"
	EMPIRE_MESSAGE_TABLE = TABLE_PREFIX + "empire_message"
	EMPIRE_NEWS_TABLE    = TABLE_PREFIX + "empire_news"
	HISTORY_CLAN_TABLE   = TABLE_PREFIX + "history_clan"
	HISTORY_EMPIRE_TABLE = TABLE_PREFIX + "history_empire"
	HISTORY_ROUND_TABLE  = TABLE_PREFIX + "history_round"
	LOCK_TABLE           = TABLE_PREFIX + "locks"
	LOG_TABLE            = TABLE_PREFIX + "log"
	LOTTERY_TABLE        = TABLE_PREFIX + "lottery"
	MARKET_TABLE         = TABLE_PREFIX + "market"
	PERMISSION_TABLE     = TABLE_PREFIX + "permission"
	SESSION_TABLE        = TABLE_PREFIX + "session"
	TURNLOG_TABLE        = TABLE_PREFIX + "turnlog"
	USER_TABLE           = TABLE_PREFIX + "users"
	VAR_TABLE            = TABLE_PREFIX + "var"
	VAR_ADJUST_TABLE     = TABLE_PREFIX + "var_adjust"
)

func (p *PHP) includes_constants_php() error {
	if !p.globals.IN_GAME {
		p.die("Access denied")
	}

	// Lookup table for translating table name token to actual table name (for setup.php)
	p.globals.tables = map[string]string{
		"{CLAN}":           CLAN_TABLE,
		"{CLAN_INVITE}":    CLAN_INVITE_TABLE,
		"{CLAN_MESSAGE}":   CLAN_MESSAGE_TABLE,
		"{CLAN_NEWS}":      CLAN_NEWS_TABLE,
		"{CLAN_RELATION}":  CLAN_RELATION_TABLE,
		"{CLAN_TOPIC}":     CLAN_TOPIC_TABLE,
		"{EMPIRE}":         EMPIRE_TABLE,
		"{EMPIRE_EFFECT}":  EMPIRE_EFFECT_TABLE,
		"{EMPIRE_MESSAGE}": EMPIRE_MESSAGE_TABLE,
		"{EMPIRE_NEWS}":    EMPIRE_NEWS_TABLE,
		"{HISTORY_CLAN}":   HISTORY_CLAN_TABLE,
		"{HISTORY_EMPIRE}": HISTORY_EMPIRE_TABLE,
		"{HISTORY_ROUND}":  HISTORY_ROUND_TABLE,
		"{LOCK}":           LOCK_TABLE,
		"{LOG}":            LOG_TABLE,
		"{LOTTERY}":        LOTTERY_TABLE,
		"{MARKET}":         MARKET_TABLE,
		"{PERMISSION}":     PERMISSION_TABLE,
		"{SESSION}":        SESSION_TABLE,
		"{TURNLOG}":        TURNLOG_TABLE,
		"{USER}":           USER_TABLE,
		"{VAR}":            VAR_TABLE,
		"{VAR_ADJUST}":     VAR_ADJUST_TABLE,
	}

	// Lookup table for translating table name to sequence name (where applicable)
	p.globals.sequences = map[string]string{
		CLAN_TABLE:           CLAN_TABLE + "_seq",
		CLAN_INVITE_TABLE:    CLAN_INVITE_TABLE + "_seq",
		CLAN_MESSAGE_TABLE:   CLAN_MESSAGE_TABLE + "_seq",
		CLAN_NEWS_TABLE:      CLAN_NEWS_TABLE + "_seq",
		CLAN_RELATION_TABLE:  CLAN_RELATION_TABLE + "_seq",
		CLAN_TOPIC_TABLE:     CLAN_TOPIC_TABLE + "_seq",
		EMPIRE_TABLE:         EMPIRE_TABLE + "_seq",
		EMPIRE_MESSAGE_TABLE: EMPIRE_MESSAGE_TABLE + "_seq",
		EMPIRE_NEWS_TABLE:    EMPIRE_NEWS_TABLE + "_seq",
		HISTORY_ROUND_TABLE:  HISTORY_ROUND_TABLE + "_seq",
		LOG_TABLE:            LOG_TABLE + "_seq",
		MARKET_TABLE:         MARKET_TABLE + "_seq",
		PERMISSION_TABLE:     PERMISSION_TABLE + "_seq",
		TURNLOG_TABLE:        TURNLOG_TABLE + "_seq",
		USER_TABLE:           USER_TABLE + "_seq",
	}

	// World variables that must be defined in order for the game to run
	p.globals.required_vars = []string{
		"lotto_current_jackpot", "lotto_yesterday_jackpot", "lotto_last_picked", "lotto_last_winner", "lotto_jackpot_increase",
		"round_time_begin", "round_time_closing", "round_time_end",
		"turns_next", "turns_next_hourly", "turns_next_daily",
	}

	// For the scope of one script execution, this is constant
	p.globals.CUR_TIME = time.Now().UTC()

	// Configurable time zones
	p.globals.timezones = map[int]string{
		-43200: "UTC-12",
		-39600: "UTC-11",
		-36000: "UTC-10",
		-34200: "UTC-9:30",
		-32400: "UTC-9",
		-28800: "UTC-8",
		-25200: "UTC-7",
		-21600: "UTC-6",
		-18000: "UTC-5",
		-14400: "UTC-4",
		-12600: "UTC-3:30",
		-10800: "UTC-3",
		-7200:  "UTC-2",
		-3600:  "UTC-1",
		0:      "UTC",
		3600:   "UTC+1",
		7200:   "UTC+2",
		10800:  "UTC+3",
		12600:  "UTC+3:30",
		14400:  "UTC+4",
		16200:  "UTC+4:30",
		18000:  "UTC+5",
		19800:  "UTC+5:30",
		20700:  "UTC+5:45",
		21600:  "UTC+6",
		23400:  "UTC+6:30",
		25200:  "UTC+7",
		28800:  "UTC+8",
		31500:  "UTC+8:45",
		32400:  "UTC+9",
		34200:  "UTC+9:30",
		36000:  "UTC+10",
		37800:  "UTC+10:30",
		39600:  "UTC+11",
		41400:  "UTC+11:30",
		43200:  "UTC+12",
		45900:  "UTC+12:45",
		46800:  "UTC+13",
		50400:  "UTC+14",
	}

	return nil
}
