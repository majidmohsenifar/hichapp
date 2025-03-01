CREATE TABLE polls (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(1024) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
