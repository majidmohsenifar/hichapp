CREATE TABLE poll_tags (
    poll_id BIGINT REFERENCES polls(id),
    tag_id BIGINT REFERENCES tags(id),
    PRIMARY KEY (poll_id, tag_id)
);
