--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- name: UserCreate :one
INSERT INTO users (u_username, u_email, u_createdate, u_lastdate)
VALUES (?, ?, datetime('now'), datetime('now'))
RETURNING u_id, u_createdate, u_lastdate;

-- name: AuthenticatedUserFetch :one
SELECT u_id, u_username, u_flags, u_comment
FROM users
WHERE u_username = ?
  AND u_password = ?;

-- name: AuthenticatedEmailFetch :one
SELECT u_id, u_username, u_flags, u_comment
FROM users
WHERE u_email = ?
  AND u_password = ?;

-- name: UserFetch :one
SELECT u_username, u_flags, u_comment, u_timezone
FROM users
WHERE u_id = ?;

-- name: UserAccessUpdate :one
UPDATE users
SET u_lastip   = ?,
    u_lastdate = datetime('now')
WHERE u_id = ?
RETURNING u_lastdate;

-- name: UserPasswordUpdate :one
UPDATE users
SET u_password = ?,
    u_lastdate = datetime('now')
WHERE u_id = ?
RETURNING u_lastdate;

-- name: UserAttributesUpdate :one
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
RETURNING u_lastdate;

-- name: WorldVarsInitialize :exec
INSERT INTO world_vars(lotto_current_jackpot, lotto_yesterday_jackpot, lotto_last_picked, lotto_last_winner,
                       lotto_jackpot_increase, round_time_begin, round_time_closing, round_time_end, turns_next,
                       turns_next_hourly, turns_next_daily)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: WorldVarsFetch :one
SELECT wv_id,
       lotto_current_jackpot,
       lotto_yesterday_jackpot,
       lotto_last_picked,
       lotto_last_winner,
       lotto_jackpot_increase,
       round_time_begin,
       round_time_closing,
       round_time_end,
       turns_next,
       turns_next_hourly,
       turns_next_daily
FROM world_vars;

-- name: WorldVarsUpdate :exec
UPDATE world_vars
SET lotto_current_jackpot   = ?,
    lotto_yesterday_jackpot = ?,
    lotto_last_picked       = ?,
    lotto_last_winner       = ?,
    lotto_jackpot_increase  = ?,
    round_time_begin        = ?,
    round_time_closing      = ?,
    round_time_end          = ?,
    turns_next              = ?,
    turns_next_hourly       = ?,
    turns_next_daily        = ?;

-- name: EmpireCreate :one
INSERT INTO empire (u_id, e_name, e_race)
VALUES (?, ?, ?)
RETURNING e_id;

-- name: EmpireActiveUserCount :one
SELECT COUNT(*)
FROM empire
WHERE u_id != 0;

-- name: UserActiveEmpires :many
SELECT e_id, e_flags
FROM empire
WHERE u_id = ?
  AND e_flags & ? = 0
ORDER BY e_id;

-- name: EmpireAttributesUpdate :exec
UPDATE empire
SET e_flags       = ?,
    e_valcode     = ?,
    e_reason      = ?,
    e_vacation    = ?,
    e_idle        = ?,
    e_era         = ?,
    e_rank        = ?,
    e_sharing     = ?,
    e_attacks     = ?,
    e_offsucc     = ?,
    e_offtotal    = ?,
    e_defsucc     = ?,
    e_deftotal    = ?,
    e_kills       = ?,
    e_score       = ?,
    e_killedby    = ?,
    e_killclan    = ?,
    e_turns       = ?,
    e_storedturns = ?,
    e_turnsused   = ?,
    e_networth    = ?,
    e_cash        = ?,
    e_food        = ?,
    e_peasants    = ?,
    e_trparm      = ?,
    e_trplnd      = ?,
    e_trpfly      = ?,
    e_trpsea      = ?,
    e_trpwiz      = ?,
    e_health      = ?,
    e_runes       = ?,
    e_indarm      = ?,
    e_indlnd      = ?,
    e_indfly      = ?,
    e_indsea      = ?,
    e_land        = ?,
    e_bldpop      = ?,
    e_bldcash     = ?,
    e_bldtrp      = ?,
    e_bldcost     = ?,
    e_bldwiz      = ?,
    e_bldfood     = ?,
    e_blddef      = ?,
    e_freeland    = ?,
    e_tax         = ?,
    e_bank        = ?,
    e_loan        = ?,
    e_mktarm      = ?,
    e_mktlnd      = ?,
    e_mktfly      = ?,
    e_mktsea      = ?,
    e_mktfood     = ?,
    e_mktperarm   = ?,
    e_mktperlnd   = ?,
    e_mktperfly   = ?,
    e_mktpersea   = ?
WHERE e_id = ?;

-- name: EmpireUpdateFlags :exec
UPDATE empire
SET e_flags = ?
WHERE e_id = ?;

-- name: ClanFetch :one
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
WHERE c_id = ?;
