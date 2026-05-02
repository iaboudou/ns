PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    content TEXT,
    image_url TEXT,
    privacy TEXT NOT NULL DEFAULT 'public',
    allowed_users TEXT,
    group_id TEXT DEFAULT NULL,
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);