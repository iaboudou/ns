
CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL, -- follow_request / follow_accepted / message / group_invite / event
    ref_id TEXT DEFAULT NULL,
    is_read BOOLEAN DEFAULT 0,
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
