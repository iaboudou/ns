
CREATE TABLE IF NOT EXISTS reactions (
    id TEXT PRIMARY KEY,
    post_or_comm_id TEXT NOT NULL,
    post_or_comm TEXT NOT NULL, -- POST / COMMENT
    type INTEGER NOT NULL DEFAULT 1, -- 1 = like
    user_id TEXT NOT NULL,
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);