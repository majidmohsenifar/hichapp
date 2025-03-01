-- name: CreatePoll :one
INSERT INTO polls (
    title
    ) VALUES (
    $1
) RETURNING *;


-- name: GetPollByID :one
SELECT * FROM polls WHERE id = $1;

-- name: GetPaginatedPollsByUserID :many
SELECT polls.* FROM polls  
LEFT JOIN votes ON polls.id = votes.poll_id AND votes.user_id = $1
WHERE votes.id IS NULL ORDER BY polls.id ASC LIMIT $2 OFFSET $3::BIGINT;

-- name: GetPaginatedPollsByUserIDTagID :many
SELECT polls.* FROM polls 
JOIN poll_tags ON polls.id = poll_tags.poll_id AND poll_tags.tag_id = $2
LEFT JOIN votes ON polls.id = votes.poll_id AND votes.user_id = $1
WHERE votes.id IS NULL ORDER BY polls.id ASC LIMIT $3 OFFSET $4::BIGINT;

-- name: GetLastCreatedPoll :one
SELECT * FROM polls ORDER BY id DESC LIMIT 1;
