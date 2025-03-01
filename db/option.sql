-- name: CreateOption :copyfrom
INSERT INTO options (
    poll_id,
    content
    ) VALUES (
    $1, $2
);

-- name: GetOptionsByPollID :many
SELECT * FROM options WHERE poll_id = $1 ORDER BY id ASC;


-- name: IncrementOptionVoteCount :exec
UPDATE options SET counts = counts + 1 WHERE id = $1;


-- name: GetOptionsByPollIDs :many
SELECT * FROM options 
WHERE poll_id = ANY($1::bigint[]);




-- name: GetOptionsContentAndCountByPollID :many
SELECT id, content, counts FROM options WHERE poll_id = $1;
