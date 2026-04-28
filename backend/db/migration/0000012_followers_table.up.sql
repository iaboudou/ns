PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS followers (
    id TEXT PRIMARY KEY,
    follower_id TEXT NOT NULL,
    following_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- pending / accepted
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(following_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(follower_id, following_id)
);