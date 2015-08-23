/*
* Database definition script for nfl_app
*
* Author: Kyle Ames
* Last Updated: August 23, 2015
*/

CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL UNIQUE,
    admin boolean NOT NULL DEFAULT FALSE,
    last_login timestamp,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS pvs (
    id integer PRIMARY KEY,
    type varchar(1) NOT NULL UNIQUE,
    seven integer NOT NULL,
    five integer NOT NULL,
    three integer NOT NULL,
    one integer NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id integer PRIMARY KEY,
    city varchar(64) NOT NULL,
    nickname varchar(64) NOT NULL,
    stadium varchar(64) NOT NULL,
    abbreviation varchar(4) NOT NULL
);

CREATE TABLE IF NOT EXISTS years (
    id integer PRIMARY KEY,
    year integer NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS weeks (
    id integer PRIMARY KEY,
    year_id integer REFERENCES years(id) ON DELETE CASCADE,
    pvs_id integer REFERENCES pvs(id),
    week integer NOT NULL,
    week_start integer
);

CREATE TABLE IF NOT EXISTS games (
    id integer PRIMARY KEY,
    week_id integer REFERENCES weeks(id) ON DELETE CASCADE,
    date integer NOT NULL,
    home_id integer REFERENCES teams(id),
    away_id integer REFERENCES teams(id),
    home_score integer DEFAULT -1,
    away_score integer DEFAULT -1
);

CREATE TABLE IF NOT EXISTS picks (
    id integer PRIMARY KEY,
    user_id integer REFERENCES users(id),
    game_id integer REFERENCES games(id),
    selection integer DEFAULT -1,
    points integer DEFAULT 0, 
    correct boolean
);

CREATE TABLE IF NOT EXISTS statistics (
    id integer PRIMARY KEY,
    user_id integer REFERENCES users(id),
    week_id integer REFERENCES weeks(id),
    zero integer,
    one integer,
    three integer,
    five integer,
    seven integer,
    winner boolean,
    lowest boolean
);
