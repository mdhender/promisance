package main

import "time"

const (
	// php/config.php

	VACATION_START = 12 * time.Hour // Delay before empire is protected
	VACATION_LIMIT = 72 * time.Hour // Minimum vacation length (not including start delay)

	// php/classes/prom_clan.php

	RELATION_INBOUND  = 2
	RELATION_OUTBOUND = 1
	RELATION_BOTH     = RELATION_OUTBOUND | RELATION_INBOUND

	// php/classes/prom_empire.php

	TURNS_TROUBLE_CASH = 1
	TURNS_TROUBLE_FOOD = 4
	TURNS_TROUBLE_LOAN = 2

	// php/classes/prom_empire_effects.php

	EMPIRE_EFFECT_PERM = "p_"
	EMPIRE_EFFECT_TIME = "m_"
	EMPIRE_EFFECT_TURN = "r_"

	// php/classes/prom_era.php

	ERA_FUTURE  = 3
	ERA_PAST    = 1
	ERA_PRESENT = 2

	// php/classes/prom_session.php

	SESSION_COOKIE = "prom_session"

	// php/classes/prom_turns.php

	TURNS_NEED_DAY    = 0x04
	TURNS_NEED_FINAL  = 0x08
	TURNS_NEED_HOUR   = 0x02
	TURNS_NEED_NORMAL = 0x01

	// php/includes/guide.php

	GUIDE_ADMIN   = 7
	GUIDE_FINANCE = 4
	GUIDE_FOREIGN = 5
	GUIDE_INFO    = 2
	GUIDE_INTRO   = 1
	GUIDE_MANAGE  = 6
	GUIDE_TURNS   = 3
	IN_GUIDE      = TRUE

	// php/includes/news.php

	CLANNEWS_MEMBER_CREATE          = 100 // e1:founder
	CLANNEWS_MEMBER_DEAD            = 104 // e1:dead member
	CLANNEWS_MEMBER_INVITE_PERM     = 108 // e1:inviter, e2:recipient
	CLANNEWS_MEMBER_INVITE_TEMP     = 107 // e1:inviter, e2:recipient
	CLANNEWS_MEMBER_JOIN            = 101 // e1:new member
	CLANNEWS_MEMBER_LEAVE           = 102 // e1:former member
	CLANNEWS_MEMBER_REMOVE          = 103 // e1:former member, e2:kicker
	CLANNEWS_MEMBER_SHARE           = 105 // e1:member
	CLANNEWS_MEMBER_UNINVITE_PERM   = 110 // e1:uninviter, e2:recipient
	CLANNEWS_MEMBER_UNINVITE_TEMP   = 109 // e1:uninviter, e2:recipient
	CLANNEWS_MEMBER_UNSHARE         = 106 // e1:member
	CLANNEWS_PERM_ASSISTANT_INHERIT = 206 // e1:new leader, e2:old leader
	CLANNEWS_PERM_GRANT_ASSISTANT   = 202 // e1:new asst, e2:grantor
	CLANNEWS_PERM_GRANT_LEADER      = 200 // e1:new leader, e2:grantor
	CLANNEWS_PERM_GRANT_MINISTER    = 204 // e1:new fa, e2:grantor
	CLANNEWS_PERM_MEMBER_INHERIT    = 208 // e1:new leader, e2:old leader
	CLANNEWS_PERM_MINISTER_INHERIT  = 207 // e1:new leader, e2:old leader
	CLANNEWS_PERM_REVOKE_ASSISTANT  = 203 // e1:old asst, e2:grantor
	CLANNEWS_PERM_REVOKE_LEADER     = 201 // e1:old leader, e2:grantor
	CLANNEWS_PERM_REVOKE_MINISTER   = 205 // e1:old fa, e2:grantor
	CLANNEWS_PROP_CHANGE_LOGO       = 303 // e1:changer
	CLANNEWS_PROP_CHANGE_PASSWORD   = 300 // e1:changer
	CLANNEWS_PROP_CHANGE_TITLE      = 301 // e1:changer
	CLANNEWS_PROP_CHANGE_URL        = 302 // e1:changer
	CLANNEWS_RECV_ALLY_DECLINE      = 410 // e2:declarer, c2:ally
	CLANNEWS_RECV_ALLY_GONE         = 411 // c2:ally
	CLANNEWS_RECV_ALLY_REQUEST      = 406 // e2:declarer, c2:ally
	CLANNEWS_RECV_ALLY_RETRACT      = 409 // e2:declarer, c2:ally
	CLANNEWS_RECV_ALLY_START        = 407 // e2:declarer, c2:ally
	CLANNEWS_RECV_ALLY_STOP         = 408 // e2:declarer, c2:ally
	CLANNEWS_RECV_WAR_GONE          = 405 // c2:opponent
	CLANNEWS_RECV_WAR_REJECT        = 404 // e2:declarer, c2:opponent
	CLANNEWS_RECV_WAR_REQUEST       = 401 // e2:declarer, c2:opponent
	CLANNEWS_RECV_WAR_RETRACT       = 403 // e2:declarer, c2:opponent
	CLANNEWS_RECV_WAR_START         = 400 // e2:declarer, c2:opponent
	CLANNEWS_RECV_WAR_STOP          = 402 // e2:declarer, c2:opponent
	CLANNEWS_SEND_ALLY_DECLINE      = 509 // e1:declarer, c2:ally
	CLANNEWS_SEND_ALLY_REQUEST      = 505 // e1:declarer, c2:ally
	CLANNEWS_SEND_ALLY_RETRACT      = 508 // e1:declarer, c2:ally
	CLANNEWS_SEND_ALLY_START        = 506 // e1:declarer, c2:ally
	CLANNEWS_SEND_ALLY_STOP         = 507 // e1:declarer, c2:ally
	CLANNEWS_SEND_WAR_REJECT        = 504 // e1:declarer, c2:opponent
	CLANNEWS_SEND_WAR_REQUEST       = 501 // e1:declarer, c2:opponent
	CLANNEWS_SEND_WAR_RETRACT       = 503 // e1:declarer, c2:opponent
	CLANNEWS_SEND_WAR_START         = 500 // e1:declarer, c2:opponent
	CLANNEWS_SEND_WAR_STOP          = 502 // e1:declarer, c2:opponent
	EMPNEWS_ATTACH_AID_RETURN       = 104 // 0:intended to return, 1:actually returned
	EMPNEWS_ATTACH_AID_SEND         = 103 // 0:convoy, 1:trparm, 2:trplnd, 3:trpfly, 4:trpsea, 5:cash, 6:runes, 7:food
	EMPNEWS_ATTACH_AID_SENDCLAN     = 105 // 0:convoy, 1:trparm, 2:trplnd, 3:trpfly, 4:trpsea, 5:cash, 6:runes, 7:food
	EMPNEWS_ATTACH_FIRST            = 100 // First attachment event, MUST be equal to the event below
	EMPNEWS_ATTACH_LAST             = 105 // Last attachment event, MUST be equal to the event above
	EMPNEWS_ATTACH_LOTTERY          = 101 // 0:winnings
	EMPNEWS_ATTACH_MARKET_RETURN    = 102 // 0:type, 1:amount, 2:price, 3:returned
	EMPNEWS_ATTACH_MARKET_SELL      = 100 // 0:type, 1:amount, 2:paid, 3:earned (minus tax)
	EMPNEWS_CLAN_ALLY_DECLINE       = 421 // no arguments
	EMPNEWS_CLAN_ALLY_GONE          = 422 // no arguments
	EMPNEWS_CLAN_ALLY_REQUEST       = 417 // no arguments
	EMPNEWS_CLAN_ALLY_RETRACT       = 420 // no arguments
	EMPNEWS_CLAN_ALLY_START         = 418 // no arguments
	EMPNEWS_CLAN_ALLY_STOP          = 419 // no arguments
	EMPNEWS_CLAN_CREATE             = 400 // no arguments
	EMPNEWS_CLAN_DISBAND            = 401 // no arguments
	EMPNEWS_CLAN_GRANT_ASSISTANT    = 407 // no arguments
	EMPNEWS_CLAN_GRANT_LEADER       = 405 // no arguments
	EMPNEWS_CLAN_GRANT_MINISTER     = 409 // no arguments
	EMPNEWS_CLAN_INHERIT_LEADER     = 406 // no arguments
	EMPNEWS_CLAN_INVITE_DISBANDED   = 428 // no arguments
	EMPNEWS_CLAN_INVITE_PERM        = 425 // no arguments
	EMPNEWS_CLAN_INVITE_TEMP        = 424 // no arguments
	EMPNEWS_CLAN_JOIN               = 402 // no arguments
	EMPNEWS_CLAN_LEAVE              = 403 // no arguments
	EMPNEWS_CLAN_REMOVE             = 404 // no arguments
	EMPNEWS_CLAN_REVOKE_ASSISTANT   = 408 // no arguments
	EMPNEWS_CLAN_REVOKE_LEADER      = 423 // no arguments
	EMPNEWS_CLAN_REVOKE_MINISTER    = 410 // no arguments
	EMPNEWS_CLAN_UNINVITE_PERM      = 427 // no arguments
	EMPNEWS_CLAN_UNINVITE_TEMP      = 426 // no arguments
	EMPNEWS_CLAN_WAR_GONE           = 416 // no arguments
	EMPNEWS_CLAN_WAR_REJECT         = 415 // no arguments
	EMPNEWS_CLAN_WAR_REQUEST        = 412 // no arguments
	EMPNEWS_CLAN_WAR_RETRACT        = 414 // no arguments
	EMPNEWS_CLAN_WAR_START          = 411 // no arguments
	EMPNEWS_CLAN_WAR_STOP           = 413 // no arguments
	EMPNEWS_MAGIC_ADVANCE           = 213 // unused
	EMPNEWS_MAGIC_BLAST             = 202 // 0:result
	EMPNEWS_MAGIC_CASH              = 208 // unused
	EMPNEWS_MAGIC_FIGHT             = 211 // 0:result (>0 = acres taken), 1:target trpwiz loss, 2:attacker trpwiz loss
	EMPNEWS_MAGIC_FOOD              = 207 // unused
	EMPNEWS_MAGIC_GATE              = 209 // unused
	EMPNEWS_MAGIC_REGRESS           = 214 // unused
	EMPNEWS_MAGIC_RUNES             = 205 // 0:result, 1:runes
	EMPNEWS_MAGIC_SHIELD            = 203 // unused
	EMPNEWS_MAGIC_SPY               = 201 // 0:result
	EMPNEWS_MAGIC_STEAL             = 212 // 0:result, 1:cash
	EMPNEWS_MAGIC_STORM             = 204 // 0:result, 1:food, 2:cash
	EMPNEWS_MAGIC_STRUCT            = 206 // 0:result, 1:buildings
	EMPNEWS_MAGIC_UNGATE            = 210 // unused
	EMPNEWS_MILITARY_AID            = 300 // 0:empire protected
	EMPNEWS_MILITARY_ARM            = 304 // 0:acres, 1:target trparm loss, 2:attacker trparm loss
	EMPNEWS_MILITARY_FLY            = 306 // 0:acres, 1:target trpfly loss, 2:attacker trpfly loss
	EMPNEWS_MILITARY_KILL           = 301 // no arguments
	EMPNEWS_MILITARY_LND            = 305 // 0:acres, 1:target trplnd loss, 2:attacker trplnd loss
	EMPNEWS_MILITARY_SEA            = 307 // 0:acres, 1:target trpsea loss, 2:attacker trpsea loss
	EMPNEWS_MILITARY_STANDARD       = 302 // 0:acres, 1:target trparm loss, 2:target trplnd loss, 3:target trpfly loss, 4:target trpsea loss,
	EMPNEWS_MILITARY_SURPRISE       = 303 // 0:acres, 1:target trparm loss, 2:target trplnd loss, 3:target trpfly loss, 4:target trpsea loss,
	SPELLRESULT_NOEFFECT            = 0
	SPELLRESULT_SHIELDED            = 1
	SPELLRESULT_SUCCESS             = 2
)

// PHPLoggingConstants (https://www.php.net/manual/en/errorfunc.constants.php)
type PHPLoggingConstants int

const (
	E_ERROR             PHPLoggingConstants = 0x0001 // Fatal run-time errors.
	E_WARNING                               = 0x0002 // Run-time warnings (non-fatal errors).
	E_PARSE                                 = 0x0004 // Compile-time parse errors.
	E_NOTICE                                = 0x0008 // Run-time notices.
	E_CORE_ERROR                            = 0x0010 // Fatal errors that occur during PHP's initial startup.
	E_CORE_WARNING                          = 0x0020 // Warnings (non-fatal errors) that occur during PHP's initial startup.
	E_COMPILE_ERROR                         = 0x0040 // Fatal compile-time errors.
	E_COMPILE_WARNING                       = 0x0080 // Compile-time warnings (non-fatal errors).
	E_USER_ERROR                            = 0x0100 // User-generated error message.
	E_USER_WARNING                          = 0x0200 // User-generated warning message.
	E_USER_NOTICE                           = 0x0400 // User-generated notice message.
	E_STRICT                                = 0x0800 // Enable to have PHP suggest changes to your code which will ensure the best interoperability and forward compatibility of your code.
	E_RECOVERABLE_ERROR                     = 0x1000 // Catchable fatal error.
	E_DEPRECATED                            = 0x2000 // Run-time notices.
	E_USER_DEPRECATED                       = 0x4000 // User-generated warning message.
	E_ALL                                   = 0xFFFF // All errors, warnings, and notices.

)

// global variable constants
// (constant per instance, not per server)
var (
	//// php/includes/constants.php
	//CUR_TIME = time.Now()

	// php/install/setup.php
	// php/util/checklang.php
	// php/util/createlocks.php
	// php/util/fixids.php
	// php/util/fixranks.php
	// php/util/worldvars.php

	IN_SCRIPT bool

	// php/index.php
	// php/install/setup.php
	// php/util/checklang.php
	// php/util/createlocks.php
	// php/util/fixids.php
	// php/util/fixranks.php
	// php/util/worldvars.php

	// php/index.php

	ROUND_CLOSING  = FALSE
	ROUND_FINISHED = TRUE
	ROUND_SIGNUP   = TRUE
	ROUND_STARTED  bool
	TXT_TIMENOTICE string
)
