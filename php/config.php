<?php
/* QM Promisance - Turn-based strategy game
 * Copyright (C) QMT Productions
 *
 * $Id: config.php 1984 2014-10-01 15:41:08Z quietust $
 */

if (!defined('IN_GAME'))
	die("Access denied");

define('DB_TYPE', 'mysql');	// Database server type (see includes/database.php for allowed values)

define('DB_SOCK', '');		// Database server socket (for local connections, UNIX only - use instead of host/port)
define('DB_HOST', '');		// Database server hostname (for Windows, or for remote connections - use instead of socket)
define('DB_PORT', '');		// Database server port (blank for default)
define('DB_USER', '');		// Database server username
define('DB_PASS', '');		// Database server password
define('DB_NAME', '');		// Database schema

define('TABLE_PREFIX', '');	// Table name prefix (in case only a single database schema is available)

define('GAME_VERSION', '4.81');
define('GAME_TITLE', 'QM Promisance '. GAME_VERSION);				// server title, may customize to liking
define('URL_BASE', 'http://www.yoursite.com/promisance/');			// the path in which Promisance is installed
define('URL_HOMEPAGE', 'http://www.yoursite.com/');				// where users will be sent when they logout
define('URL_FORUMS', 'http://www.yoursite.com/forums/viewforum.php?f=N');	// where your site's forums are located
										// (comment out if you have none)
define('MAIL_ADMIN', 'admin@yoursite.com');					// administrative contact address
define('MAIL_VALIDATE', 'promisance@yoursite.com');				// source address for validation emails
define('TXT_NEWS', '<span style="color:white">Welcome to '. GAME_TITLE .'!</span>');
										// news message displayed at top of all pages
										// undefine if you don't want any
define('TXT_RULES', 'Additional information can be found in our forums.<br />'	// extra rules to display on signup page
	.'Please read through the FAQs before contacting an Administrator.');	// undefine if you don't want any
define('TXT_EMAIL', 'Be sure to check out http://www.qmtpro.com/ and tell your friends about us!');
										// custom text to display in signup email
define('DEFAULT_LANGUAGE', 'en-US');	// Default language pack to use, also used when not logged in
define('BASE_LANGUAGE', 'en-US');	// Base language pack, used for anything not defined in current/default

define('EMPIRES_PER_USER', 1);		// How many empires can be owned at once by a particular user?
define('TURNS_PROTECTION', 200);	// Duration of protection
define('TURNS_INITIAL', 100);		// Turns given on signup
define('TURNS_MAXIMUM', 250);		// Max accumulated turns
define('TURNS_STORED', 100);		// Max stored turns
define('TURNS_VALIDATE', 150);		// How long before validation is necessary
define('TURNS_ERA', 500);		// Minimum number of turns that must be spent in an era before one can advance or regress

define('VACATION_START', 12);		// Delay before empire is protected
define('VACATION_LIMIT', 72);		// Minimum vacation length (not including start delay)

define('TURNS_FREQ', 10);	// how often to give turns
define('TURNS_OFFSET', 0);	// offset (in minutes) for giving turns, relative to round start
define('TURNS_OFFSET_HOURLY', 0);	// offset (in minutes) for performing hourly events, relative to round start
define('TURNS_OFFSET_DAILY', 60*12);	// offset (in minutes) for performing daily events, relative to round start
define('TURNS_CRONTAB', TRUE);	// use "turns.php" to give out turns, scheduled via crontab; otherwise, trigger on page loads
define('TURNS_CRONLOG', TRUE);	// if TURNS_CRONTAB is disabled, store turn logs in the database for retrieval by turns.php
define('TURNS_COUNT', 1);	// how many turns to give during each period
define('TURNS_UNSTORE', 1);	// how many turns to release from Stored Turns at once

define('IDLE_TIMEOUT_NEW', 3);		// Remove new empire if idle for this many days before being prompted to validate (create and abandon)
define('IDLE_TIMEOUT_VALIDATE', 2);	// Remove new empire if prompted to validate but fails to do so within this many days (bad email address)
define('IDLE_TIMEOUT_ABANDON', 14);	// Remove established empire if idle for this many days (and not on vacation or disabled)
define('IDLE_TIMEOUT_KILLED', 2);	// Remove dead empire after this many days if never logged in to see notification
define('IDLE_TIMEOUT_DELETE', 3);	// Remove deleted empire after this many days (unless still under protection, in which case it is immediate)

define('LOTTERY_MAXTICKETS', 3);	// Maximum number of lottery tickets per empire
define('LOTTERY_JACKPOT', 1000000000);	// Base jackpot

define('BUILD_COST', 3500);	// Base building cost
define('DROP_DELAY', 60*60*12);	// Minimum delay (in seconds) between gaining land and dropping it

define('BANK_SAVERATE', 4.0);	// Base savings interest rate
define('BANK_LOANRATE', 7.5);	// Base loan interest rate

define('PUBMKT_START', 6);	// Hours before goods will arrive on public market
define('PUBMKT_MINTIME', -1);	// Number of hours before users can manually remove items (-1 to disallow)
define('PUBMKT_MAXTIME', 72);	// Number of hours before items are automatically removed (-1 to disallow)
define('PUBMKT_MINSELL', 0);	// Minimum percentage of troops, per shipment, that can be sold on public market (0-100)
define('PUBMKT_MAXSELL', 25);	// Maximum percentage of troops, total, that can be sold on public market (0-100)
define('PUBMKT_MINFOOD', 0);	// Same as MINSELL, except for food
define('PUBMKT_MAXFOOD', 90);	// Same as MAXSELL, except for food

define('CLAN_ENABLE', TRUE);	// Master enable for clans
define('CLAN_MINJOIN', 72);	// Empires can't leave clans until they've been a member for this many hours
define('CLAN_MINREJOIN', 24);	// Empires can't create/join a new clan until this many hours after they left
define('CLAN_MINSHARE', 2);	// Unsharing clan forces takes this many hours to take effect
define('CLAN_MINRELATE', 48);	// How long a clan must wait before it can modify an alliance or war slot
define('CLAN_MAXALLY', 3);	// Maximum number of alliances (outbound and inbound)
define('CLAN_MAXWAR', 3);	// Maximum number of wars (outbound only)
define('CLAN_VIEW_STAT', TRUE);	// Allow clan leaders to view the detailed Empire Status of members
define('CLAN_VIEW_AID', TRUE);	// Display summary of all foreign aid sent between clan members on Clan Management page
define('CLAN_INVITE_TIME', 48);	// How long clan invites last before they expire
define('CLAN_LATE_JOIN', TRUE);	// Allow empires to join clans during the cooldown period

define('PVTM_MAXSELL', 8000);	// Percentage of troops that can be sold on private market (0-10000)
define('PVTM_SHOPBONUS', 0.20);	// Percentage of private market cost bonus for which shops are responsible (0-1)
define('PVTM_TRPARM', 500);	// Base market costs for each unit
define('PVTM_TRPLND', 1000);
define('PVTM_TRPFLY', 2000);
define('PVTM_TRPSEA', 3000);
define('PVTM_FOOD', 30);

define('INDUSTRY_MULT', 2.5);	// Industry output multiplier
define('MAX_ATTACKS', 30);	// Maximum number of attacks

define('AID_ENABLE', TRUE);	// Enable sending foreign aid
define('AID_MAXCREDITS', 5);	// Maximum number of aid credits that can be accumulated at once
define('AID_DELAY', 60 * 60);	// Once you send too much aid, how many seconds before you can send more
define('MESSAGES_MAXCREDITS', 5);	// Maximum number of new private messages that can be sent at once
define('MESSAGES_DELAY', 10 * 60);	// Once you send too many messages, how many seconds before you can send more
define('FRIEND_MAGIC_ENABLE', FALSE);	// Enable casting spells on friendly empires
define('SCORE_ENABLE', FALSE);	// Enable keeping score for empires attacking each other
define('MAGIC_ALLOW_REGRESS', FALSE);	// Enables usage of the "Regress to Previous Era" spell
define('GRAVEYARD_DISCLOSE', FALSE);	// Reveal user account name of empires in the Graveyard

define('CLANSTATS_MINSIZE', 1);	// Minimum member count for inclusion on Top Clans page
define('TOPEMPIRES_COUNT', 50);	// Maximum number of empires to list on Top Empires page
define('TOPPLAYERS_COUNT', 50);	// Maximum number of user accounts to list on All Time Top Players page

define('SIGNUP_CLOSED_USER', FALSE);	// Don't allow creation of new user accounts from signup page
define('SIGNUP_CLOSED_EMPIRE', FALSE);	// Don't allow creation of new empires from signup page

define('VALIDATE_REQUIRE', TRUE);	// Require users to validate their empires before they can continue playing
define('VALIDATE_ALLOW', TRUE);	// Allow users to validate their own empires
define('VALIDATE_RESEND', 60 * 60);	// How long users must wait between resending their validation code

define('COUNTER_TEMPLATE', 'counter2.png');	// Digit style for "registered empires" counter
						// Set to empty string to use plain bold text
define('COUNTER_DIGITS', 3);	// Number of digits for "registered empires" counter
				// Set to 0 to determine automatically

define('BONUS_TURNS', TRUE);	// Allow users to collect 1 hour worth of bonus turns each day.
				// Turns are collected by clicking a banner (defined below) or,
				// if none are defined, a "Bonus Turns" button at the top of the page.

$banners = array();		// Define banners below. Do NOT use paid advertisements if bonus
				// turns are enabled, as it provides an incentive for click fraud.
// $banners[] = array('label' => 'Hover text', 'image' => 'Image URL', 'url' => 'Click URL', 'width' => '468', 'height' => '60', 'ismap' => '1' if imagemap, '0' if not);

// Stylesheets for use in-game
$styles = array(
	'qmt' => array('file' => 'qmt.css', 'name' => 'QMT Black'),
	'ezclan' => array('file' => 'ezclan.css', 'name' => 'EZClan Blue'),
);
define('DEFAULT_STYLE', 'qmt');	// Default style, also used when not logged in
define('DEFAULT_TIMEZONE', -21600);	// Default timezone for new accounts
					// Specify UTC offset in seconds, must be defined in constants.php
define('DEFAULT_DATEFORMAT', 'r');	// Default date/time format for new accounts

// Default values for newly created empires
$empire_defaults = array(
	// Resources
	'e_cash'	=> 100000,
	'e_food'	=> 10000,
	'e_runes'	=> 0,

	// Units
	'e_peasants'	=> 500,
	'e_trparm'	=> 100,
	'e_trplnd'	=> 50,
	'e_trpfly'	=> 20,
	'e_trpsea'	=> 10,
	'e_trpwiz'	=> 0,

	// Buildings
	'e_land'	=> 250,	// acre counts below MUST add up to this value!
	'e_bldpop'	=> 20,
	'e_bldcash'	=> 10,
	'e_bldtrp'	=> 0,
	'e_bldcost'	=> 5,
	'e_bldwiz'	=> 0,
	'e_bldfood'	=> 15,
	'e_blddef'	=> 0,
	'e_freeland'	=> 200,

	// Private Market supplies
	'e_mktarm'	=> 4000,
	'e_mktlnd'	=> 3000,
	'e_mktfly'	=> 2000,
	'e_mktsea'	=> 1000,
	'e_mktfood'	=> 100000,

	// Others
	'e_indarm'	=> 25,
	'e_indlnd'	=> 25,
	'e_indfly'	=> 25,
	'e_indsea'	=> 25,
	'e_health'	=> 100,
	'e_tax'		=> 10,
);

define('LOG_ENABLE', FALSE);	// Enable logging of all in-game events
				// Warnings are logged regardless of this setting
define('DEBUG_FOOTER', FALSE);	// Display page generation time, memory usage, and query count at bottom of page
?>
