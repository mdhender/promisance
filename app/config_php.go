// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

const (
	DB_TYPE = "mysql" // Database server type (see includes/database.php for allowed values)

	DB_SOCK = "" // Database server socket (for local connections, UNIX only - use instead of host/port)
	DB_HOST = "" // Database server hostname (for Windows, or for remote connections - use instead of socket)
	DB_PORT = "" // Database server port (blank for default)
	DB_USER = "" // Database server username
	DB_PASS = "" // Database server password
	DB_NAME = "" // Database schema

	TABLE_PREFIX = "" // Table name prefix (in case only a single database schema is available)

	GAME_VERSION = "4.81"
	GAME_TITLE   = "QM Promisance " + GAME_VERSION        // server title, may customize to liking
	URL_BASE     = "https://www.yoursite.com/promisance/" // the path in which Promisance is installed
	URL_HOMEPAGE = "https://www.yoursite.com/"            // where users will be sent when they logout

	// where your site's forums are located
	// (comment out if you have none)
	URL_FORUMS = "https://www.yoursite.com/forums/viewforum.php?f=N"
	MAIL_ADMIN = "admin@yoursite.com" // administrative contact address

	// source address for validation emails
	// news message displayed at top of all pages
	// undefine if you don"t want any
	MAIL_VALIDATE = "promisance@yoursite.com"
	TXT_NEWS      = "<span style='color:white'>Welcome to " + GAME_TITLE + "!</span>"

	// extra rules to display on signup page
	// undefine if you don't want any
	TXT_RULES = "Additional information can be found in our forums.<br />Please read through the FAQs before contacting an Administrator."

	// custom text to display in signup email
	TXT_EMAIL = "Be sure to check out https://www.qmtpro.com/ and tell your friends about us!"

	DEFAULT_LANGUAGE = "en-US" // Default language pack to use, also used when not logged in
	BASE_LANGUAGE    = "en-US" // Base language pack, used for anything not defined in current/default

	EMPIRES_PER_USER = 1   // How many empires can be owned at once by a particular user?
	TURNS_PROTECTION = 200 // Duration of protection
	TURNS_INITIAL    = 100 // Turns given on signup
	TURNS_MAXIMUM    = 250 // Max accumulated turns
	TURNS_STORED     = 100 // Max stored turns
	TURNS_VALIDATE   = 150 // How long before validation is necessary
	TURNS_ERA        = 500 // Minimum number of turns that must be spent in an era before one can advance or regress

	TURNS_FREQ          = 10      // how often to give turns
	TURNS_OFFSET        = 0       // offset (in minutes) for giving turns, relative to round start
	TURNS_OFFSET_HOURLY = 0       // offset (in minutes) for performing hourly events, relative to round start
	TURNS_OFFSET_DAILY  = 12 * 60 // offset (in minutes) for performing daily events, relative to round start
	TURNS_CRONTAB       = TRUE    // use "turns.php" to give out turns, scheduled via crontab; otherwise, trigger on page loads
	TURNS_CRONLOG       = TRUE    // if TURNS_CRONTAB is disabled, store turn logs in the database for retrieval by turns.php
	TURNS_COUNT         = 1       // how many turns to give during each period
	TURNS_UNSTORE       = 1       // how many turns to release from Stored Turns at once

	IDLE_TIMEOUT_NEW      = 3  // Remove new empire if idle for this many days before being prompted to validate (create and abandon)
	IDLE_TIMEOUT_VALIDATE = 2  // Remove new empire if prompted to validate but fails to do so within this many days (bad email address)
	IDLE_TIMEOUT_ABANDON  = 14 // Remove established empire if idle for this many days (and not on vacation or disabled)
	IDLE_TIMEOUT_KILLED   = 2  // Remove dead empire after this many days if never logged in to see notification
	IDLE_TIMEOUT_DELETE   = 3  // Remove deleted empire after this many days (unless still under protection, in which case it is immediate)

	LOTTERY_MAXTICKETS = 3             // Maximum number of lottery tickets per empire
	LOTTERY_JACKPOT    = 1_000_000_000 // Base jackpot

	BUILD_COST = 3_500        // Base building cost
	DROP_DELAY = 12 * 60 * 60 // Minimum delay (in seconds) between gaining land and dropping it

	BANK_SAVERATE = 4.0 // Base savings interest rate
	BANK_LOANRATE = 7.5 // Base loan interest rate

	PUBMKT_START   = 6  // Hours before goods will arrive on public market
	PUBMKT_MINTIME = -1 // Number of hours before users can manually remove items (-1 to disallow)
	PUBMKT_MAXTIME = 72 // Number of hours before items are automatically removed (-1 to disallow)
	PUBMKT_MINSELL = 0  // Minimum percentage of troops, per shipment, that can be sold on public market (0-100)
	PUBMKT_MAXSELL = 25 // Maximum percentage of troops, total, that can be sold on public market (0-100)
	PUBMKT_MINFOOD = 0  // Same as MINSELL, except for food
	PUBMKT_MAXFOOD = 90 // Same as MAXSELL, except for food

	CLAN_ENABLE      = TRUE // Master enable for clans
	CLAN_MINJOIN     = 72   // Empires can"t leave clans until they"ve been a member for this many hours
	CLAN_MINREJOIN   = 24   // Empires can"t create/join a new clan until this many hours after they left
	CLAN_MINSHARE    = 2    // Unsharing clan forces takes this many hours to take effect
	CLAN_MINRELATE   = 48   // How long a clan must wait before it can modify an alliance or war slot
	CLAN_MAXALLY     = 3    // Maximum number of alliances (outbound and inbound)
	CLAN_MAXWAR      = 3    // Maximum number of wars (outbound only)
	CLAN_VIEW_STAT   = TRUE // Allow clan leaders to view the detailed Empire Status of members
	CLAN_VIEW_AID    = TRUE // Display summary of all foreign aid sent between clan members on Clan Management page
	CLAN_INVITE_TIME = 48   // How long clan invites last before they expire
	CLAN_LATE_JOIN   = TRUE // Allow empires to join clans during the cooldown period

	PVTM_MAXSELL   = 8000 // Percentage of troops that can be sold on private market (0-10000)
	PVTM_SHOPBONUS = 0.20 // Percentage of private market cost bonus for which shops are responsible (0-1)
	PVTM_TRPARM    = 500  // Base market costs for each unit
	PVTM_TRPLND    = 1000
	PVTM_TRPFLY    = 2000
	PVTM_TRPSEA    = 3000
	PVTM_FOOD      = 30

	INDUSTRY_MULT = 2.5 // Industry output multiplier
	MAX_ATTACKS   = 30  // Maximum number of attacks

	AID_ENABLE     = TRUE    // Enable sending foreign aid
	AID_MAXCREDITS = 5       // Maximum number of aid credits that can be accumulated at once
	AID_DELAY      = 60 * 60 // Once you send too much aid, how many seconds before you can send more

	MESSAGES_MAXCREDITS = 5       // Maximum number of new private messages that can be sent at once
	MESSAGES_DELAY      = 10 * 60 // Once you send too many messages, how many seconds before you can send more

	FRIEND_MAGIC_ENABLE = FALSE // Enable casting spells on friendly empires
	SCORE_ENABLE        = FALSE // Enable keeping score for empires attacking each other

	MAGIC_ALLOW_REGRESS = FALSE // Enables usage of the "Regress to Previous Era" spell
	GRAVEYARD_DISCLOSE  = FALSE // Reveal user account name of empires in the Graveyard

	CLANSTATS_MINSIZE = 1  // Minimum member count for inclusion on Top Clans page
	TOPEMPIRES_COUNT  = 50 // Maximum number of empires to list on Top Empires page
	TOPPLAYERS_COUNT  = 50 // Maximum number of user accounts to list on All Time Top Players page

	SIGNUP_CLOSED_USER   = FALSE // Don"t allow creation of new user accounts from signup page
	SIGNUP_CLOSED_EMPIRE = FALSE // Don"t allow creation of new empires from signup page

	VALIDATE_REQUIRE = TRUE    // Require users to validate their empires before they can continue playing
	VALIDATE_ALLOW   = TRUE    // Allow users to validate their own empires
	VALIDATE_RESEND  = 60 * 60 // How long users must wait between resending their validation code

	// Digit style for "registered empires" counter
	// Set to empty string to use plain bold text
	COUNTER_TEMPLATE = "counter2.png"

	// Number of digits for "registered empires" counter
	// Set to 0 to determine automatically
	COUNTER_DIGITS = 3

	// Allow users to collect 1 hour worth of bonus turns each day.
	// Turns are collected by clicking a banner (defined below) or,
	// if none are defined, a "Bonus Turns" button at the top of the page.
	BONUS_TURNS = TRUE

	DEFAULT_DATEFORMAT = "r"   // Default date/time format for new accounts
	DEFAULT_STYLE      = "qmt" // Default style, also used when not logged in

	// Default timezone for new accounts
	// Specify UTC offset in seconds, must be defined in constants.php
	DEFAULT_TIMEZONE = -21600

	// Enable logging of all in-game events
	LOG_ENABLE = FALSE
	// Warnings are logged regardless of this setting
	// Display page generation time, memory usage, and query count at bottom of page
	DEBUG_FOOTER = FALSE
)

func (p *PHP) config_php() error {
	if !p.constants.IN_GAME {
		p.die("Access denied")
	}

	// Define banners below. Do NOT use paid advertisements if bonus
	// turns are enabled, as it provides an incentive for click fraud.
	// $banners[] = array('label' => 'Hover text', 'image' => 'Image URL', 'url' => 'Click URL', 'width' => '468', 'height' => '60', 'ismap' => '1' if imagemap, '0' if not);

	// Stylesheets for use in-game
	p.globals.styles = map[string]css_file_t{
		"qmt":    {file: "qmt.css", name: "QMT Black"},
		"ezclan": {file: "ezclan.css", name: "EZClan Blue"},
	}

	// Default values for newly created empires
	p.globals.empire_defaults = map[string]int{
		// Resources
		"e_cash":  100000,
		"e_food":  10000,
		"e_runes": 0,

		// Units
		"e_peasants": 500,
		"e_trparm":   100,
		"e_trplnd":   50,
		"e_trpfly":   20,
		"e_trpsea":   10,
		"e_trpwiz":   0,

		// Buildings
		"e_land":     250, // acre counts below MUST add up to this value!
		"e_bldpop":   20,
		"e_bldcash":  10,
		"e_bldtrp":   0,
		"e_bldcost":  5,
		"e_bldwiz":   0,
		"e_bldfood":  15,
		"e_blddef":   0,
		"e_freeland": 200,

		// Private Market supplies
		"e_mktarm":  4000,
		"e_mktlnd":  3000,
		"e_mktfly":  2000,
		"e_mktsea":  1000,
		"e_mktfood": 100000,

		// Others
		"e_indarm": 25,
		"e_indlnd": 25,
		"e_indfly": 25,
		"e_indsea": 25,
		"e_health": 100,
		"e_tax":    10,
	}

	return nil
}
