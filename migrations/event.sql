CREATE TABLE events (
    uuid UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    event_type VARCHAR(100) NOT NULL,
    version INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    amount INTEGER,
    oldpoint INTEGER,
    newpoint INTEGER,

    UNIQUE (user_id, version)
);