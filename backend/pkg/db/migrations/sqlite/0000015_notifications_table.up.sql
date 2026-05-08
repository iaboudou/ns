PRAGMA foreign_keys = ON;

CREATE TABLE notifications (
    id TEXT PRIMARY KEY,
    sender_id TEXT NOT NULL,
    type TEXT NOT NULL,
    group_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES users(id)
);