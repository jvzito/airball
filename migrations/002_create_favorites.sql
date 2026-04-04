CREATE TABLE IF NOT EXISTS favorites (
    id         SERIAL PRIMARY KEY,
    user_id    INT REFERENCES users(id) ON DELETE CASCADE,
    player_id  INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, player_id)
);
