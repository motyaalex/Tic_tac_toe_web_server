-- PostgreSQL schema for AP1-Go-T04B
-- Create the database separately if it does not exist, then execute this file in that database.

CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY,
    login         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS games (
    id           UUID PRIMARY KEY,
    board        JSONB NOT NULL,
    player_x     UUID,
    player_o     UUID,
    state        TEXT NOT NULL,
    turn_player  UUID,
    winner       UUID,
    vs_computer  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_games_state ON games (state);
