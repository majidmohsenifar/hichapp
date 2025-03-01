-- name: CreateVote :one
INSERT INTO votes (
    poll_id,
    option_id,
    user_id,
    is_skipped
    ) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: GetVoteByPollIDAndUserID :one
SELECT * FROM votes WHERE poll_id = $1 AND user_id = $2;
