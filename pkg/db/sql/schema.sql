CREATE TABLE IF NOT EXISTS
  users (
    id TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    active BOOLEAN NOT NULL DEFAULT TRUE
  );

CREATE TABLE IF NOT EXISTS
  sessions (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users (id),
    expiry INTEGER NOT NULL
  );

CREATE TABLE IF NOT EXISTS
  connections (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users (id),
    provider TEXT NOT NULL,
    email TEXT NOT NULL,
    token TEXT NOT NULL,
    refresh_token TEXT
  )