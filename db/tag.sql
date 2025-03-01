-- name: CreateTag :copyfrom
INSERT INTO tags (
    name
    ) VALUES (
    $1
);

-- name: GetTagsByPollIDs :many
SELECT * FROM tags JOIN poll_tags ON tags.id = poll_tags.tag_id 
WHERE poll_id = ANY($1::bigint[]);


-- name: GetTagByName :one
SELECT * FROM tags WHERE name = $1;
