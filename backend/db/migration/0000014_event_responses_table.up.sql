PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS event_responses (
    event_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at DATE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (event_id, user_id),

    FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);