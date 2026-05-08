PRAGMA foreign_keys = ON;

CREATE TABLE notification_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    notification_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    is_read INTEGER DEFAULT 0,
    FOREIGN KEY(notification_id) REFERENCES notifications(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);