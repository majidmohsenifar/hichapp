CREATE TABLE votes (
    id BIGSERIAL PRIMARY KEY,
    poll_id BIGINT NOT NULL REFERENCES polls(id),
    option_id BIGINT REFERENCES options(id),
    user_id BIGINT NOT NULL,
    is_skipped BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
)
