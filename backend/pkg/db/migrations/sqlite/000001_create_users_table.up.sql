PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    nickname TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    birthday TEXT NOT NULL,
    gender TEXT NOT NULL,
    profile_image TEXT NOT NULL,
    last_notif_seen DATE DEFAULT CURRENT_TIMESTAMP,
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    about_me TEXT DEFAULT NULL,
    account_privacy BOOLEAN DEFAULT 0
);