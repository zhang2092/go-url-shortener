-- name: CreateUserUrl :one
INSERT INTO user_relate_url (
  user_id, short_url, origin_url, status, expire_at
) VALUES (
  $1, $2, $3, 0, $4
)
RETURNING *;

-- name: ListUrlByUser :many
SELECT *
FROM user_relate_url
WHERE user_id = $1
ORDER BY id DESC;