/*
* Database definition script for nfl_app
*
* Author: Kyle Ames
* Last Updated: August 17, 2014
*/

/* Clean up existing database if needed */
CREATE DATABASE IF NOT EXISTS nfl WITH OWNER nfl_app;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL,
    admin boolean NOT NULL DEFAULT FALSE,
    last_login timestamp,
    password text NOT NULL
);

CREATE TABLE IF NOT EXISTS pvs (
    id SERIAL PRIMARY KEY,
    type varchar(1) NOT NULL,
    seven integer NOT NULL,
    five integer NOT NULL,
    three integer NOT NULL,
    one integer NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    city varchar(64) NOT NULL,
    nickname varchar(64) NOT NULL,
    stadium varchar(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS years (
    id SERIAL PRIMARY KEY,
    year integer NOT NULL
);

CREATE TABLE IF NOT EXISTS weeks (
    id SERIAL PRIMARY KEY,
    year_id integer REFERENCES years ON DELETE CASCADE,
    pvs_id integer REFERENCES pvs
);

CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    week_id integer REFERENCES weeks ON DELETE CASCADE,
    date timestamp NOT NULL,
    home_id integer REFERENCES teams,
    away_id integer REFERENCES teams,
    home_score integer DEFAULT -1,
    away_score integer DEFAULT -1
);

CREATE TABLE IF NOT EXISTS picks (
    id SERIAL PRIMARY KEY,
    user_id integer REFERENCES users,
    game_id integer REFERENCES games,
    selection integer DEFAULT -1,
    points integer DEFAULT 0, 
    correct boolean
);
