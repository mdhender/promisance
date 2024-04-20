-- QM Promisance - Turn-based strategy game
-- Copyright (C) QMT Productions
--
-- $Id: prom.sqlite 1983 2014-10-01 15:18:43Z quietust $

DROP TABLE IF EXISTS clan;
CREATE TABLE clan
(
    c_id        INTEGER PRIMARY KEY,
    c_name      TEXT    NOT NULL,           -- varchar(8)   NOT NULL DEFAULT '',
    c_password  TEXT    NOT NULL,           -- varchar(255) NOT NULL DEFAULT '',
    c_members   INTEGER NOT NULL DEFAULT 0, -- tinyint      NOT NULL DEFAULT 0,
    e_id_leader INTEGER NOT NULL DEFAULT 0, -- int unsigned NOT NULL DEFAULT 0,
    e_id_asst   INTEGER NOT NULL DEFAULT 0, -- int unsigned NOT NULL DEFAULT 0,
    e_id_fa1    INTEGER NOT NULL DEFAULT 0, -- int unsigned NOT NULL DEFAULT 0,
    e_id_fa2    INTEGER NOT NULL DEFAULT 0, -- int unsigned NOT NULL DEFAULT 0,
    c_title     TEXT    NOT NULL,           -- varchar(255) NOT NULL DEFAULT '',
    c_url       TEXT    NOT NULL,           -- varchar(255) NOT NULL DEFAULT '',
    c_pic       TEXT    NOT NULL            -- varchar(255) NOT NULL DEFAULT ''
);
CREATE INDEX clan_c_name ON clan (c_name);
CREATE INDEX clan_c_members ON clan (c_members);

DROP TABLE IF EXISTS clan_invite;
CREATE TABLE clan_invite
(
    ci_id    INTEGER PRIMARY KEY,
    c_id     int unsigned     NOT NULL DEFAULT 0,
    e_id_1   int unsigned     NOT NULL DEFAULT 0,
    e_id_2   int unsigned     NOT NULL DEFAULT 0,
    ci_flags tinyint unsigned NOT NULL DEFAULT 0,
    ci_time  int              NOT NULL DEFAULT 0
);
CREATE INDEX clan_invite_c_id ON clan_invite (c_id);
CREATE INDEX clan_invite_e_id_2 ON clan_invite (e_id_2);

DROP TABLE IF EXISTS clan_message;
CREATE TABLE clan_message
(
    cm_id    INTEGER PRIMARY KEY,
    ct_id    int unsigned     NOT NULL DEFAULT 0,
    e_id     int unsigned     NOT NULL DEFAULT 0,
    cm_body  text             NOT NULL,
    cm_time  int              NOT NULL DEFAULT 0,
    cm_flags tinyint unsigned NOT NULL DEFAULT 0
);
CREATE INDEX clan_message_e_id ON clan_message (e_id);
CREATE INDEX clan_message_cm_time ON clan_message (cm_time);

DROP TABLE IF EXISTS clan_news;
CREATE TABLE clan_news
(
    cn_id    INTEGER PRIMARY KEY,
    cn_time  int               NOT NULL DEFAULT 0,
    c_id     int unsigned      NOT NULL DEFAULT 0,
    e_id_1   int unsigned      NOT NULL DEFAULT 0,
    c_id_2   int unsigned      NOT NULL DEFAULT 0,
    e_id_2   int unsigned      NOT NULL DEFAULT 0,
    cn_event smallint unsigned NOT NULL DEFAULT 0
);
CREATE INDEX clan_news_c_id ON clan_news (c_id);
CREATE INDEX clan_news_e_id_1 ON clan_news (e_id_1);
CREATE INDEX clan_news_c_id_2 ON clan_news (c_id_2);
CREATE INDEX clan_news_e_id_2 ON clan_news (e_id_2);
CREATE INDEX clan_news_cn_event ON clan_news (cn_event);

DROP TABLE IF EXISTS clan_relation;
CREATE TABLE clan_relation
(
    cr_id    INTEGER PRIMARY KEY,
    c_id_1   int unsigned     NOT NULL DEFAULT 0,
    c_id_2   int unsigned     NOT NULL DEFAULT 0,
    cr_flags tinyint unsigned NOT NULL DEFAULT 0,
    cr_time  int              NOT NULL DEFAULT 0
);
CREATE INDEX clan_relation_c_id_1 ON clan_relation (c_id_1);
CREATE INDEX clan_relation_c_id_2 ON clan_relation (c_id_2);
CREATE INDEX clan_relation_cr_flags ON clan_relation (cr_flags);

DROP TABLE IF EXISTS clan_topic;
CREATE TABLE clan_topic
(
    ct_id      INTEGER PRIMARY KEY,
    c_id       int unsigned     NOT NULL DEFAULT 0,
    ct_subject varchar(255)     NOT NULL DEFAULT '',
    ct_flags   tinyint unsigned NOT NULL DEFAULT 0
);
CREATE INDEX clan_topic_c_id ON clan_topic (c_id);

DROP TABLE IF EXISTS empire;
CREATE TABLE empire
(
    e_id          INTEGER PRIMARY KEY,
    u_id          INTEGER NOT NULL, -- int unsigned       NOT NULL DEFAULT 0,
    u_oldid       INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_signupdate  TIMESTAMP,        -- int                NOT NULL DEFAULT 0,
    e_flags       INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_valcode     TEXT,             -- varchar(32)        NOT NULL DEFAULT '',
    e_reason      TEXT,             -- varchar(255)       NOT NULL DEFAULT '',
    e_vacation    INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_idle        INTEGER,          -- int                NOT NULL DEFAULT 0,
    e_name        TEXT    NOT NULL, -- varchar(255)       NOT NULL DEFAULT '',
    e_race        INTEGER NOT NULL, -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_era         INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_rank        INTEGER,          -- mediumint unsigned NOT NULL DEFAULT 0,
    c_id          INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    c_oldid       INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_sharing     INTEGER,          -- tinyint            NOT NULL DEFAULT 0,
    e_attacks     INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_offsucc     INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_offtotal    INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_defsucc     INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_deftotal    INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_kills       INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_score       INTEGER,          -- int                NOT NULL DEFAULT 0,
    e_killedby    INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_killclan    INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_turns       INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_storedturns INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_turnsused   INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_networth    INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_cash        INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_food        INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_peasants    INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_trparm      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_trplnd      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_trpfly      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_trpsea      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_trpwiz      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_health      INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_runes       INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_indarm      INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_indlnd      INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_indfly      INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_indsea      INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_land        INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_bldpop      INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_bldcash     INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_bldtrp      INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_bldcost     INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_bldwiz      INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_bldfood     INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_blddef      INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_freeland    INTEGER,          -- int unsigned       NOT NULL DEFAULT 0,
    e_tax         INTEGER,          -- tinyint unsigned   NOT NULL DEFAULT 0,
    e_bank        INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_loan        INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_mktarm      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_mktlnd      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_mktfly      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_mktsea      INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_mktfood     INTEGER,          -- bigint unsigned    NOT NULL DEFAULT 0,
    e_mktperarm   INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_mktperlnd   INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_mktperfly   INTEGER,          -- smallint unsigned  NOT NULL DEFAULT 0,
    e_mktpersea   INTEGER           -- smallint unsigned  NOT NULL DEFAULT 0
);
CREATE INDEX empire_u_id ON empire (u_id);
CREATE INDEX empire_u_oldid ON empire (u_oldid);
CREATE INDEX empire_e_flags ON empire (e_flags);
CREATE INDEX empire_c_id ON empire (c_id);

DROP TABLE IF EXISTS empire_effect;
CREATE TABLE empire_effect
(
    e_id     int unsigned   NOT NULL DEFAULT 0,
    ef_name  varbinary(255) NOT NULL DEFAULT '',
    ef_value int            NOT NULL DEFAULT 0,
    PRIMARY KEY (e_id, ef_name)
);

DROP TABLE IF EXISTS empire_message;
CREATE TABLE empire_message
(
    m_id      INTEGER PRIMARY KEY,
    m_id_ref  int unsigned     NOT NULL DEFAULT 0,
    m_time    int              NOT NULL DEFAULT 0,
    e_id_src  int unsigned     NOT NULL DEFAULT 0,
    e_id_dst  int unsigned     NOT NULL DEFAULT 0,
    m_subject varchar(255)     NOT NULL DEFAULT '',
    m_body    text             NOT NULL,
    m_flags   tinyint unsigned NOT NULL DEFAULT 0
);
CREATE INDEX empire_message_m_time ON empire_message (m_time);
CREATE INDEX empire_message_e_id_src ON empire_message (e_id_src);
CREATE INDEX empire_message_e_id_dst ON empire_message (e_id_dst);
CREATE INDEX empire_message_m_flags ON empire_message (m_flags);

DROP TABLE IF EXISTS empire_news;
CREATE TABLE empire_news
(
    n_id     INTEGER PRIMARY KEY,
    n_time   int               NOT NULL DEFAULT 0,
    e_id_src int unsigned      NOT NULL DEFAULT 0,
    c_id_src int unsigned      NOT NULL DEFAULT 0,
    e_id_dst int unsigned      NOT NULL DEFAULT 0,
    c_id_dst int unsigned      NOT NULL DEFAULT 0,
    n_event  smallint unsigned NOT NULL DEFAULT 0,
    n_d0     bigint            NOT NULL DEFAULT 0,
    n_d1     bigint            NOT NULL DEFAULT 0,
    n_d2     bigint            NOT NULL DEFAULT 0,
    n_d3     bigint            NOT NULL DEFAULT 0,
    n_d4     bigint            NOT NULL DEFAULT 0,
    n_d5     bigint            NOT NULL DEFAULT 0,
    n_d6     bigint            NOT NULL DEFAULT 0,
    n_d7     bigint            NOT NULL DEFAULT 0,
    n_d8     bigint            NOT NULL DEFAULT 0,
    n_flags  tinyint unsigned  NOT NULL DEFAULT 0
);
CREATE INDEX empire_news_e_id_src ON empire_news (e_id_src);
CREATE INDEX empire_news_c_id_src ON empire_news (c_id_src);
CREATE INDEX empire_news_e_id_dst ON empire_news (e_id_dst);
CREATE INDEX empire_news_c_id_dst ON empire_news (c_id_dst);
CREATE INDEX empire_news_n_event ON empire_news (n_event);
CREATE INDEX empire_news_n_flags ON empire_news (n_flags);

DROP TABLE IF EXISTS history_clan;
CREATE TABLE history_clan
(
    hr_id       smallint        NOT NULL DEFAULT 0,
    hc_id       int unsigned    NOT NULL DEFAULT 0,
    hc_members  smallint        NOT NULL DEFAULT 0,
    hc_name     varchar(8)      NOT NULL DEFAULT '',
    hc_title    varchar(255)    NOT NULL DEFAULT '',
    hc_totalnet bigint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (hr_id, hc_id)
);

DROP TABLE IF EXISTS history_empire;
CREATE TABLE history_empire
(
    hr_id       smallint           NOT NULL DEFAULT 0,
    he_flags    tinyint unsigned   NOT NULL DEFAULT 0,
    u_id        int unsigned       NOT NULL DEFAULT 0,
    he_id       int unsigned       NOT NULL DEFAULT 0,
    he_name     varchar(255)       NOT NULL DEFAULT '',
    he_race     varchar(64)        NOT NULL DEFAULT '',
    he_era      varchar(64)        NOT NULL DEFAULT '',
    hc_id       int unsigned       NOT NULL DEFAULT 0,
    he_offsucc  smallint unsigned  NOT NULL DEFAULT 0,
    he_offtotal smallint unsigned  NOT NULL DEFAULT 0,
    he_defsucc  smallint unsigned  NOT NULL DEFAULT 0,
    he_deftotal smallint unsigned  NOT NULL DEFAULT 0,
    he_kills    smallint unsigned  NOT NULL DEFAULT 0,
    he_score    int                NOT NULL DEFAULT 0,
    he_networth bigint unsigned    NOT NULL DEFAULT 0,
    he_land     int unsigned       NOT NULL DEFAULT 0,
    he_rank     mediumint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY (hr_id, he_id)
);

DROP TABLE IF EXISTS history_round;
CREATE TABLE history_round
(
    hr_id             INTEGER PRIMARY KEY,
    hr_name           varchar(64)        NOT NULL DEFAULT '',
    hr_description    text               NOT NULL,
    hr_startdate      int                NOT NULL DEFAULT 0,
    hr_stopdate       int                NOT NULL DEFAULT 0,
    hr_flags          tinyint unsigned   NOT NULL DEFAULT 0,
    hr_smallclansize  tinyint unsigned   NOT NULL DEFAULT 0,
    hr_smallclans     smallint unsigned  NOT NULL DEFAULT 0,
    hr_allclans       smallint unsigned  NOT NULL DEFAULT 0,
    hr_nonclanempires mediumint unsigned NOT NULL DEFAULT 0,
    hr_liveempires    mediumint unsigned NOT NULL DEFAULT 0,
    hr_deadempires    mediumint unsigned NOT NULL DEFAULT 0,
    hr_delempires     mediumint unsigned NOT NULL DEFAULT 0,
    hr_allempires     mediumint unsigned NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS locks;
CREATE TABLE locks
(
    lock_id bigint unsigned NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS log;
CREATE TABLE log
(
    log_id     INTEGER PRIMARY KEY,
    log_time   int unsigned NOT NULL DEFAULT 0,
    log_type   int unsigned NOT NULL DEFAULT 0,
    log_ip     varchar(40)  NOT NULL DEFAULT '',
    log_page   varchar(32)  NOT NULL DEFAULT '',
    log_action varchar(64)  NOT NULL DEFAULT '',
    log_locks  varchar(64)  NOT NULL DEFAULT '',
    log_text   text         NOT NULL,
    u_id       int unsigned NOT NULL DEFAULT 0,
    e_id       int unsigned NOT NULL DEFAULT 0,
    c_id       int unsigned NOT NULL DEFAULT 0
);

DROP TABLE IF EXISTS lottery;
CREATE TABLE lottery
(
    e_id     int unsigned    NOT NULL DEFAULT 0,
    l_ticket int unsigned    NOT NULL DEFAULT 0,
    l_cash   bigint unsigned NOT NULL DEFAULT 0
);
CREATE INDEX lottery_e_id ON lottery (e_id);
CREATE INDEX lottery_l_ticket ON lottery (l_ticket);

DROP TABLE IF EXISTS market;
CREATE TABLE market
(
    k_id    INTEGER PRIMARY KEY,
    k_type  tinyint unsigned NOT NULL DEFAULT 0,
    e_id    int unsigned     NOT NULL DEFAULT 0,
    k_amt   bigint unsigned  NOT NULL DEFAULT 0,
    k_price int unsigned     NOT NULL DEFAULT 0,
    k_time  int              NOT NULL DEFAULT 0
);
CREATE INDEX market_e_id ON market (e_id);
CREATE INDEX market_k_type ON market (k_type);
CREATE INDEX market_k_time ON market (k_time);

DROP TABLE IF EXISTS permission;
CREATE TABLE permission
(
    p_id         INTEGER PRIMARY KEY,
    p_type       tinyint unsigned NOT NULL DEFAULT 0,
    p_criteria   varchar(255)     NOT NULL DEFAULT '',
    p_comment    varchar(255)     NOT NULL DEFAULT '',
    p_reason     varchar(255)     NOT NULL DEFAULT '',
    p_createtime int unsigned     NOT NULL DEFAULT 0,
    p_updatetime int unsigned     NOT NULL DEFAULT 0,
    p_lasthit    int unsigned     NOT NULL DEFAULT 0,
    p_hitcount   int unsigned     NOT NULL DEFAULT 0,
    p_expire     int unsigned     NOT NULL DEFAULT 0
);
CREATE INDEX permission_p_type ON permission (p_type);
CREATE INDEX permission_p_expire ON permission (p_expire);

DROP TABLE IF EXISTS session;
CREATE TABLE session
(
    sess_id   varbinary(64) PRIMARY KEY,
    sess_time int unsigned NOT NULL DEFAULT 0,
    sess_data blob         NOT NULL
);
CREATE INDEX session_sess_time ON session (sess_time);

DROP TABLE IF EXISTS turnlog;
CREATE TABLE turnlog
(
    turn_id       INTEGER PRIMARY KEY,
    turn_time     int unsigned     NOT NULL DEFAULT 0,
    turn_ticks    int unsigned     NOT NULL DEFAULT 0,
    turn_interval int unsigned     NOT NULL DEFAULT 0,
    turn_type     tinyint unsigned NOT NULL DEFAULT 0,
    turn_text     text             NOT NULL
);
CREATE INDEX turnlog_turn_type ON turnlog (turn_type);

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    u_id         INTEGER PRIMARY KEY,
    u_username   TEXT NOT NULL, -- varchar(255)     NOT NULL DEFAULT '',
    u_password   TEXT,          -- varchar(255)     NOT NULL DEFAULT '',
    u_flags      INTEGER,       -- tinyint unsigned NOT NULL DEFAULT 0,
    u_name       TEXT,          -- varchar(255)     NOT NULL DEFAULT '',
    u_email      TEXT NOT NULL, -- varchar(255)     NOT NULL DEFAULT '',
    u_comment    TEXT,          -- varchar(255)     NOT NULL DEFAULT '',
    u_timezone   INTEGER,       -- mediumint        NOT NULL DEFAULT 0,
    u_style      TEXT,          -- varchar(32)      NOT NULL DEFAULT '',
    u_lang       TEXT,          -- varchar(16)      NOT NULL DEFAULT '',
    u_dateformat TEXT,          -- varchar(64)      NOT NULL DEFAULT '',
    u_lastip     TEXT,          -- varchar(40)      NOT NULL DEFAULT '',
    u_kills      INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_deaths     INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_offsucc    INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_offtotal   INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_defsucc    INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_deftotal   INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_numplays   INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_sucplays   INTEGER,       -- int unsigned     NOT NULL DEFAULT 0,
    u_avgrank    REAL,          -- double           NOT NULL DEFAULT 0,
    u_bestrank   REAL,          -- double           NOT NULL DEFAULT 0,
    u_createdate TIMESTAMP,     -- int              NOT NULL DEFAULT 0,
    u_lastdate   TIMESTAMP,     -- int              NOT NULL DEFAULT 0,
    UNIQUE (u_username),
    UNIQUE (u_email)
);
CREATE INDEX users_u_flags ON users (u_flags);

DROP TABLE IF EXISTS var;
CREATE TABLE var
(
    v_name  TEXT NOT NULL, -- varbinary(255) NOT NULL DEFAULT '',
    v_value TEXT NOT NULL, -- varbinary(255) NOT NULL DEFAULT '',
    PRIMARY KEY (v_name)
);

DROP TABLE IF EXISTS var_adjust;
CREATE TABLE var_adjust
(
    v_name   varbinary(255) NOT NULL DEFAULT '',
    v_offset bigint         NOT NULL DEFAULT 0
);
