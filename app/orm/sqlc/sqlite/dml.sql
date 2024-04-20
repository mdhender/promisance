--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

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

-- name: UserCreate :one
INSERT INTO users (u_username, u_email, u_createdate, u_lastdate)
VALUES (?, ?, datetime('now'), datetime('now'))
RETURNING u_id, u_createdate, u_lastdate;

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

