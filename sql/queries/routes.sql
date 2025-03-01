-- name: GetRoutes :many
SELECT *
FROM routes;

-- name: GetRoute :one
SELECT *
FROM routes
WHERE id = $1
LIMIT 1;

-- name: CreateRoute :one
INSERT INTO routes (name)
VALUES ($1)
RETURNING id;

-- name: DeleteRoute :one
DELETE FROM routes
WHERE id = $1
RETURNING id;

-- name: UpdateRoute :exec
UPDATE routes
SET name = $1
WHERE id = sqlc.arg(id);
