
CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    image_url TEXT,
    created_at DATE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);