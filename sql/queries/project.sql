-- name: CreateProject :one
INSERT INTO projects (name, status, link, description, metadata, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateProjectStatus :exec
UPDATE projects
SET status = $2, updated_at = $3
WHERE name = $1;

-- name: UpdateProjectMetadata :exec
UPDATE projects
SET metadata = $2, updated_at = $3
WHERE name = $1;

-- name: GetProjectByName :one
SELECT * FROM projects
WHERE name = $1
LIMIT 1;

-- name: ListAllProjects :many
SELECT * FROM projects
ORDER BY created_at DESC;

-- name: ListNProjects :many
SELECT * FROM projects
ORDER BY created_at DESC
LIMIT $1;

-- name: ListProjectsByStatus :many
SELECT * FROM projects
WHERE status = $1
ORDER BY created_at DESC;

-- name: DeleteProjectByName :exec
DELETE FROM projects
WHERE name = $1;

-- name: GetProjectIdByName :one
SELECT id FROM projects
WHERE name = $1;

-- name: GetProjectStatusById :one
SELECT status FROM projects
WHERE id = $1;

-- name: UpdateProjectStatusById :one
UPDATE projects
SET status = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: CheckProjectExistsByName :one
SELECT EXISTS (
    SELECT 1 FROM projects WHERE name = $1
) AS exists;
